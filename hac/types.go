package hac

type ScriptType string

const (
	ScriptGroovy     ScriptType = "groovy"
	ScriptBeanShell  ScriptType = "beanshell"
	ScriptJavaScript ScriptType = "javascript"
)

type FlexQuery struct {
	FlexibleSearchQuery string `form:"flexibleSearchQuery"`
	SQLQuery            string `form:"sqlQuery"`
	User                string `form:"user" default:"admin"`
	Locale              string `form:"locale" default:"en"`
	MaxCount            int    `form:"maxCount" default:"200"`
	DataSource          string `form:"dataSource" default:"master"`
	Commit              bool   `form:"commit"`
	NoAnalyze           bool   `form:"noAnalyze"`
}

type FlexExecuteOptions struct {
	ColumnBlacklist []string
	NoBlacklist     bool
}

type FlexSearchResponse struct {
	CatalogVersionsAsString string     `json:"catalogVersionsAsString,omitempty"`
	DataSourceId            string     `json:"dataSourceId,omitempty"`
	Exception               any        `json:"exception,omitempty"`
	ExceptionStackTrace     string     `json:"exceptionStackTrace,omitempty"`
	ExecutionTime           float64    `json:"executionTime,omitempty"`
	Headers                 []string   `json:"headers,omitempty"`
	ParametersAsString      string     `json:"parametersAsString,omitempty"`
	Query                   string     `json:"query,omitempty"`
	RawExecution            bool       `json:"rawExecution,omitempty"`
	ResultCount             float64    `json:"resultCount,omitempty"`
	ResultList              [][]string `json:"resultList,omitempty"`
}

type FlexException struct {
	Message string `json:"message"`
}

type GroovyRequest struct {
	Script     string     `form:"script"`
	ScriptType ScriptType `form:"scriptType" default:"groovy"`
	Commit     bool       `form:"commit,keepzero"`
}

type GroovyResponse struct {
	Result     string `json:"executionResult"`
	Stacktrace string `json:"stacktraceText"`
	Output     string `json:"outputText"`
}

type ImportValidationEnum string

const (
	ImportStrict  ImportValidationEnum = "IMPORT_STRICT"
	ImportRelaxed ImportValidationEnum = "IMPORT_RELAXED"
)

type ExportValidationEnum string

const (
	ExportOnly            ExportValidationEnum = "EXPORT_ONLY"
	ExportReimportStrict  ExportValidationEnum = "EXPORT_REIMPORT_STRICT"
	ExportReimportRelaxed ExportValidationEnum = "EXPORT_REIMPORT_RELAXED"
)

type ImpexImportRequest struct {
	ScriptContent        string               `form:"scriptContent"`
	ValidationEnum       ImportValidationEnum `form:"validationEnum" default:"IMPORT_STRICT"`
	MaxThreads           int                  `form:"maxThreads,keepzero" default:"8"`
	Encoding             string               `form:"encoding" default:"UTF-8"`
	LegacyMode           bool                 `form:"legacyMode"`
	EnableCodeExecution  bool                 `form:"enableCodeExecution"`
	DistributedMode      bool                 `form:"distributedMode"`
	SldEnabled           bool                 `form:"sldEnabled"`
	LegacyMode_          string               `form:"_legacyMode" default:"on"`
	EnableCodeExecution_ string               `form:"_enableCodeExecution" default:"on"`
	DistributedMode_     string               `form:"_distributedMode" default:"on"`
	SldEnabled_          string               `form:"_sldEnabled" default:"on"`
}

type ImpexExportRequest struct {
	ScriptContent  string               `form:"scriptContent"`
	ValidationEnum ExportValidationEnum `form:"validationEnum" default:"EXPORT_ONLY"`
	Encoding       string               `form:"encoding" default:"UTF-8"`
}

type TypeAttributesRequest struct {
	TypeCode string `form:"type"`
}

type TypeAttributesResponse struct {
	Exists     bool     `json:"exists"`
	Attributes []string `json:"attributes"`
}

type PKAnalyzeRequest struct {
	PKString string `form:"pkString"`
}

type PKAnalyzeResponse struct {
	PKString           string   `json:"pkString"`
	CounterBased       bool     `json:"counterBased"`
	PKAsHex            string   `json:"pkAsHex"`
	PKAsBinary         string   `json:"pkAsBinary"`
	PKTypeCode         int      `json:"pkTypeCode"`
	PKClusterId        int      `json:"pkClusterId"`
	PKCreationTime     int64    `json:"pkCreationTime"`
	PKCreationDate     string   `json:"pkCreationDate"`
	Bits               []string `json:"bits"`
	PossibleException  any      `json:"possibleException"`
	PKMilliCnt         int      `json:"pkMilliCnt"`
	PKComposedTypeCode string   `json:"pkComposedTypeCode"`
}
