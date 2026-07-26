package hac

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

type ImpexService struct {
	client *HACClient
}

func NewImpexService(c *HACClient) *ImpexService {
	return &ImpexService{client: c}
}

func (s *ImpexService) Import(
	ctx context.Context,
	q ImpexImportRequest,
) (string, error) {

	q = applyDefaults(q)
	form := s.client.buildForm(q)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := s.client.doRequest(
		ctx,
		http.MethodPost,
		"console/impex/import",
		strings.NewReader(form.Encode()),
		headers,
	)
	if err != nil {
		return "", err
	}

	body, err := readAllBody(resp)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(string(body))
	if doc.Error != nil {
		return "", doc.Error
	}

	resultTag := doc.Find("span", "id", "impexResult")
	if resultTag.Error != nil {
		resultTag = doc.Find("div", "class", "impexResult")
	}

	var result string
	if resultTag.Error != nil {
		result = ""
	} else {
		result = resultTag.Attrs()["data-result"]
		if result == "" {
			result = resultTag.FullText()
		}
	}

	return result, nil
}

func (s *ImpexService) FetchTypeAndAttributes(t TypeAttributesRequest) (*TypeAttributesResponse, error) {
	form := s.client.buildForm(t)

	req, err := s.client.doRequest(
		context.Background(),
		http.MethodPost,
		"console/impex/typeAndAttributes",
		strings.NewReader(form.Encode()),
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)
	if err != nil {
		return nil, err
	}

	body, err := readAllBody(req)
	if err != nil {
		return nil, err
	}

	var out TypeAttributesResponse
	if e := json.Unmarshal([]byte(body), &out); e != nil {
		return nil, fmt.Errorf("decode failed: %w body: %s", e, body)
	}
	return &out, nil
}

func (s *ImpexService) Export(ctx context.Context, q ImpexExportRequest) (string, string, error) {
	q = applyDefaults(q)
	form := s.client.buildForm(q)

	resp, err := s.client.doRequest(
		ctx,
		http.MethodPost,
		"console/impex/export",
		strings.NewReader(form.Encode()),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	)

	if err != nil {
		return "", "", fmt.Errorf("failed to execute impex export: %w", err)
	}

	body, err := readAllBody(resp)
	if err != nil {
		return "", "", err
	}

	doc := soup.HTMLParse(string(body))
	if doc.Error != nil {
		return "", "", doc.Error
	}

	resultTag := doc.Find("span", "id", "impexResult")

	if resultTag.Error != nil {
		return "", "", fmt.Errorf("missing impexResult tag")
	}

	result := resultTag.Attrs()["data-result"]

	parent := doc.Find("div", "id", "downloadExportResultData")
	if parent.Error != nil {
		return result, "", fmt.Errorf("missing downloadExportResultData container")
	}

	link := parent.Find("a")
	if link.Error != nil {
		return result, "", fmt.Errorf("missing download link inside container")
	}

	href, ok := link.Attrs()["href"]

	if !ok || href == "" {
		return result, "", fmt.Errorf("download link missing href attribute")
	}

	return result, href, nil
}

func (s *ImpexService) DownloadExportZip(downloadPath string) ([]byte, error) {
	base := s.client.baseURL
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}

	full := base + "console/impex/" + downloadPath

	req, err := http.NewRequest(http.MethodGet, full, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", base+"console/impex/export")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", s.client.userAgent)
	if s.client.csrf != "" {
		req.Header.Set("X-CSRF-TOKEN", s.client.csrf)
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "zip") && !strings.Contains(ct, "octet-stream") {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("not zip. status %d body:\n%s", resp.StatusCode, b)
	}

	return io.ReadAll(resp.Body)
}

func (s *ImpexService) UploadZip(
	ctx context.Context,
	filename string,
	fileData []byte,
) (string, error) {

	b64 := base64.StdEncoding.EncodeToString(fileData)
	jobCode := fmt.Sprintf("importCronJob_%d", time.Now().UnixNano())

	script := fmt.Sprintf(`
import de.hybris.platform.impex.model.cronjob.ImpExImportCronJobModel
import de.hybris.platform.impex.model.cronjob.ImpExImportJobModel
import de.hybris.platform.impex.enums.ImpExValidationModeEnum
import de.hybris.platform.impex.model.ImpExMediaModel

def zipBytes = Base64.getDecoder().decode("%s")

ImpExMediaModel jobMedia = modelService.create(ImpExMediaModel.class)
jobMedia.setCode("%s_zip")
jobMedia.setCatalogVersion(null)
modelService.save(jobMedia)
mediaService.setStreamForMedia(jobMedia, new ByteArrayInputStream(zipBytes), "import.zip", "application/zip")

// ImpExMediaModel mediasMedia = modelService.create(ImpExMediaModel.class)
// mediasMedia.setCode("%s_medias")
// mediasMedia.setCatalogVersion(null)
// modelService.save(mediasMedia)
// mediaService.setStreamForMedia(mediasMedia, new ByteArrayInputStream(zipBytes), "import.zip", "application/zip")

ImpExImportJobModel importJob = flexibleSearchService
    .search("SELECT {PK} FROM {ImpExImportJob} WHERE {code} = 'ImpEx-Import'")
    .getResult()
    .get(0)

ImpExImportCronJobModel importCronJob = modelService.create(ImpExImportCronJobModel.class)
importCronJob.setCode("%s")
importCronJob.setJob(importJob)
importCronJob.setJobMedia(jobMedia)
importCronJob.setZipentry("importscript.impex")
// importCronJob.setMediasMedia(mediasMedia)
// importCronJob.setUnzipMediasMedia(true)
importCronJob.setMode(ImpExValidationModeEnum.IMPORT_STRICT)
importCronJob.setMaxThreads(10)
importCronJob.setLogToDatabase(true)
importCronJob.setDumpingAllowed(true)
modelService.save(importCronJob)

cronJobService.performCronJob(importCronJob, true)
modelService.refresh(importCronJob)

println "STATUS:" + importCronJob.status
println "RESULT:" + importCronJob.result
println "LASTLINE:" + importCronJob.lastSuccessfulLine
println "VALUECOUNT:" + importCronJob.valueCount
if (importCronJob.unresolvedDataStore) {
    println "UNRESOLVED:" + importCronJob.unresolvedDataStore.code
}
`, b64, jobCode, jobCode, jobCode)

	resp, err := s.client.Groovy.Execute(ctx, GroovyRequest{
		Script: script,
	})
	if err != nil {
		return "", fmt.Errorf("groovy execute failed: %w", err)
	}
	if resp == nil {
		return "", fmt.Errorf("empty groovy response")
	}

	output := resp.Output

	status := extractField(output, "STATUS:")
	result := extractField(output, "RESULT:")

	if result != "" && result != "SUCCESS" {
		return output, fmt.Errorf("import finished with result=%s status=%s", result, status)
	}

	return output, nil
}

func extractField(output, prefix string) string {
	for line := range strings.SplitSeq(output, "\n") {
		if after, ok :=strings.CutPrefix(line, prefix); ok  {
			return strings.TrimSpace(after)
		}
	}
	return ""
}
