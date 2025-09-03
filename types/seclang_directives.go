package types

import "fmt"

type SeclangDirective interface {
	ToSeclang() string
}

type ChainableDirective interface {
	SeclangDirective
	GetMetadata() Metadata
	GetActions() *SeclangActions
	GetTransformations() Transformations
	ToSeclangWithIdent(string) string
	GetChainedDirective() ChainableDirective
	AppendChainedDirective(ChainableDirective)
	NonDisruptiveActionsCount() int
}

type CommentDirective struct {
	Kind     Kind            `yaml:"kind"`
	Metadata CommentMetadata `yaml:",inline"`
}

func (d CommentDirective) ToSeclang() string {
	return d.Metadata.ToSeclang()
}

type ConfigurationDirective struct {
	Kind      Kind                       `yaml:"kind"`
	Metadata  *CommentMetadata           `yaml:",inline"`
	Name      ConfigurationDirectiveType `yaml:"name"`
	Parameter string                     `yaml:"parameter"`
}

type ConfigurationDirectiveType string

const (
	SecAuditLogStorageDir         ConfigurationDirectiveType = "SecAuditLogStorageDir"
	SecAuditLogDirMode            ConfigurationDirectiveType = "SecAuditLogDirMode"
	SecAuditEngine                ConfigurationDirectiveType = "SecAuditEngine"
	SecAuditLogFileMode           ConfigurationDirectiveType = "SecAuditLogFileMode"
	SecAuditLog2                  ConfigurationDirectiveType = "SecAuditLog2"
	SecAuditLog                   ConfigurationDirectiveType = "SecAuditLog"
	SecAuditLogFormat             ConfigurationDirectiveType = "SecAuditLogFormat"
	SecAuditLogParts              ConfigurationDirectiveType = "SecAuditLogParts"
	SecAuditLogRelevantStatus     ConfigurationDirectiveType = "SecAuditLogRelevantStatus"
	SecAuditLogType               ConfigurationDirectiveType = "SecAuditLogType"
	SecUploadKeepFiles            ConfigurationDirectiveType = "SecUploadKeepFiles"
	SecTmpSaveUploadedFiles       ConfigurationDirectiveType = "SecTmpSaveUploadedFiles"
	SecUploadDir                  ConfigurationDirectiveType = "SecUploadDir"
	SecUploadFileLimit            ConfigurationDirectiveType = "SecUploadFileLimit"
	SecUploadFileMode             ConfigurationDirectiveType = "SecUploadFileMode"
	SecComponentSignature         ConfigurationDirectiveType = "SecComponentSignature"
	SecServerSignature            ConfigurationDirectiveType = "SecServerSignature"
	SecWebAppId                   ConfigurationDirectiveType = "SecWebAppId"
	SecMarker                     ConfigurationDirectiveType = "SecMarker"
	SecConnEngine                 ConfigurationDirectiveType = "SecConnEngine"
	SecContentInjection           ConfigurationDirectiveType = "SecContentInjection"
	SecArgumentsLimit             ConfigurationDirectiveType = "SecArgumentsLimit"
	SecDebugLog                   ConfigurationDirectiveType = "SecDebugLog"
	SecDebugLogLevel              ConfigurationDirectiveType = "SecDebugLogLevel"
	SecGeoLookupDb                ConfigurationDirectiveType = "SecGeoLookupDb"
	SecGsbLookupDb                ConfigurationDirectiveType = "SecGsbLookupDb"
	SecPcreMatchLimit             ConfigurationDirectiveType = "SecPcreMatchLimit"
	SecPcreMatchLimitRecursion    ConfigurationDirectiveType = "SecPcreMatchLimitRecursion"
	SecRequestBodyJsonDepthLimit  ConfigurationDirectiveType = "SecRequestBodyJsonDepthLimit"
	SecRequestBodyAccess          ConfigurationDirectiveType = "SecRequestBodyAccess"
	SecRequestBodyInMemoryLimit   ConfigurationDirectiveType = "SecRequestBodyInMemoryLimit"
	SecRequestBodyLimit           ConfigurationDirectiveType = "SecRequestBodyLimit"
	SecRequestBodyLimitAction     ConfigurationDirectiveType = "SecRequestBodyLimitAction"
	SecRequestBodyNoFilesLimit    ConfigurationDirectiveType = "SecRequestBodyNoFilesLimit"
	SecResponseBodyMimeType       ConfigurationDirectiveType = "SecResponseBodyMimeType"
	SecResponseBodyMimeTypesClear ConfigurationDirectiveType = "SecResponseBodyMimeTypesClear"
	SecResponseBodyAccess         ConfigurationDirectiveType = "SecResponseBodyAccess"
	SecResponseBodyLimit          ConfigurationDirectiveType = "SecResponseBodyLimit"
	SecResponseBodyLimitAction    ConfigurationDirectiveType = "SecResponseBodyLimitAction"
	SecRuleEngine                 ConfigurationDirectiveType = "SecRuleEngine"
	SecCookieFormat               ConfigurationDirectiveType = "SecCookieFormat"
	SecCookieV0Separator          ConfigurationDirectiveType = "SecCookieV0Separator"
	SecDataDir                    ConfigurationDirectiveType = "SecDataDir"
	SecStatusEngine               ConfigurationDirectiveType = "SecStatusEngine"
	SecTmpDir                     ConfigurationDirectiveType = "SecTmpDir"
	SecUnicodeMapFile             ConfigurationDirectiveType = "SecUnicodeMapFile"
	SecArgumentSeparator          ConfigurationDirectiveType = "SecArgumentSeparator"
	SecChrootDir                  ConfigurationDirectiveType = "SecChrootDir"
	SecCollectionTimeout          ConfigurationDirectiveType = "SecCollectionTimeout"
	SecConnReadStateLimit         ConfigurationDirectiveType = "SecConnReadStateLimit"
	SecConnWriteStateLimit        ConfigurationDirectiveType = "SecConnWriteStateLimit"
	SecDisableBackendCompression  ConfigurationDirectiveType = "SecDisableBackendCompression"
	SecGuardianLog                ConfigurationDirectiveType = "SecGuardianLog"
	SecHashEngine                 ConfigurationDirectiveType = "SecHashEngine"
	SecHashKey                    ConfigurationDirectiveType = "SecHashKey"
	SecHashParam                  ConfigurationDirectiveType = "SecHashParam"
	SecHashMethodRx               ConfigurationDirectiveType = "SecHashMethodRx"
	SecHashMethodPm               ConfigurationDirectiveType = "SecHashMethodPm"
	SecHttpBlKey                  ConfigurationDirectiveType = "SecHttpBlKey"
	SecInterceptOnError           ConfigurationDirectiveType = "SecInterceptOnError"
	SecRemoteRulesFailAction      ConfigurationDirectiveType = "SecRemoteRulesFailAction"
	SecRuleInheritance            ConfigurationDirectiveType = "SecRuleInheritance"
	SecRulePerfTime               ConfigurationDirectiveType = "SecRulePerfTime"
	SecSensorId                   ConfigurationDirectiveType = "SecSensorId"
	SecStreamInBodyInspection     ConfigurationDirectiveType = "SecStreamInBodyInspection"
	SecStreamOutBodyInspection    ConfigurationDirectiveType = "SecStreamOutBodyInspection"
	SecXmlExternalEntity          ConfigurationDirectiveType = "SecXmlExternalEntity"
	SecCacheTransformations       ConfigurationDirectiveType = "SecCacheTransformations"
)

var (
	configDirectiveTypes = map[string]ConfigurationDirectiveType{
		"SecAuditLogStorageDir":         SecAuditLogStorageDir,
		"SecAuditLogDirMode":            SecAuditLogDirMode,
		"SecAuditEngine":                SecAuditEngine,
		"SecAuditLogFileMode":           SecAuditLogFileMode,
		"SecAuditLog2":                  SecAuditLog2,
		"SecAuditLog":                   SecAuditLog,
		"SecAuditLogFormat":             SecAuditLogFormat,
		"SecAuditLogParts":              SecAuditLogParts,
		"SecAuditLogRelevantStatus":     SecAuditLogRelevantStatus,
		"SecAuditLogType":               SecAuditLogType,
		"SecUploadKeepFiles":            SecUploadKeepFiles,
		"SecTmpSaveUploadedFiles":       SecTmpSaveUploadedFiles,
		"SecUploadDir":                  SecUploadDir,
		"SecUploadFileLimit":            SecUploadFileLimit,
		"SecUploadFileMode":             SecUploadFileMode,
		"SecComponentSignature":         SecComponentSignature,
		"SecServerSignature":            SecServerSignature,
		"SecWebAppId":                   SecWebAppId,
		"SecMarker":                     SecMarker,
		"SecConnEngine":                 SecConnEngine,
		"SecContentInjection":           SecContentInjection,
		"SecArgumentsLimit":             SecArgumentsLimit,
		"SecDebugLog":                   SecDebugLog,
		"SecDebugLogLevel":              SecDebugLogLevel,
		"SecGeoLookupDb":                SecGeoLookupDb,
		"SecGsbLookupDb":                SecGsbLookupDb,
		"SecPcreMatchLimit":             SecPcreMatchLimit,
		"SecPcreMatchLimitRecursion":    SecPcreMatchLimitRecursion,
		"SecRequestBodyJsonDepthLimit":  SecRequestBodyJsonDepthLimit,
		"SecRequestBodyAccess":          SecRequestBodyAccess,
		"SecRequestBodyInMemoryLimit":   SecRequestBodyInMemoryLimit,
		"SecRequestBodyLimit":           SecRequestBodyLimit,
		"SecRequestBodyLimitAction":     SecRequestBodyLimitAction,
		"SecRequestBodyNoFilesLimit":    SecRequestBodyNoFilesLimit,
		"SecResponseBodyMimeType":       SecResponseBodyMimeType,
		"SecResponseBodyMimeTypesClear": SecResponseBodyMimeTypesClear,
		"SecResponseBodyAccess":         SecResponseBodyAccess,
		"SecResponseBodyLimit":          SecResponseBodyLimit,
		"SecResponseBodyLimitAction":    SecResponseBodyLimitAction,
		"SecRuleEngine":                 SecRuleEngine,
		"SecCookieFormat":               SecCookieFormat,
		"SecCookieV0Separator":          SecCookieV0Separator,
		"SecDataDir":                    SecDataDir,
		"SecStatusEngine":               SecStatusEngine,
		"SecTmpDir":                     SecTmpDir,
		"SecUnicodeMapFile":             SecUnicodeMapFile,
		"SecArgumentSeparator":          SecArgumentSeparator,
		"SecChrootDir":                  SecChrootDir,
		"SecCollectionTimeout":          SecCollectionTimeout,
		"SecConnReadStateLimit":         SecConnReadStateLimit,
		"SecConnWriteStateLimit":        SecConnWriteStateLimit,
		"SecDisableBackendCompression":  SecDisableBackendCompression,
		"SecGuardianLog":                SecGuardianLog,
		"SecHashEngine":                 SecHashEngine,
		"SecHashKey":                    SecHashKey,
		"SecHashParam":                  SecHashParam,
		"SecHashMethodRx":               SecHashMethodRx,
		"SecHashMethodPm":               SecHashMethodPm,
		"SecHttpBlKey":                  SecHttpBlKey,
		"SecInterceptOnError":           SecInterceptOnError,
		"SecRemoteRulesFailAction":      SecRemoteRulesFailAction,
		"SecRuleInheritance":            SecRuleInheritance,
		"SecRulePerfTime":               SecRulePerfTime,
		"SecSensorId":                   SecSensorId,
		"SecStreamInBodyInspection":     SecStreamInBodyInspection,
		"SecStreamOutBodyInspection":    SecStreamOutBodyInspection,
		"SecXmlExternalEntity":          SecXmlExternalEntity,
		"SecCacheTransformations":       SecCacheTransformations,
	}
)

func NewConfigurationDirective() *ConfigurationDirective {
	c := new(ConfigurationDirective)
	c.Kind = ConfigurationKind
	c.Metadata = new(CommentMetadata)
	return c
}

func (c *ConfigurationDirective) SetName(name string) error {
	constName, ok := configDirectiveTypes[name]
	if !ok {
		return fmt.Errorf("Invalid configuration directive name: %s", name)
	}
	c.Name = constName
	return nil
}

func (c ConfigurationDirective) GetMetadata() Metadata {
	return c.Metadata
}

// TODO: add quotes around the value when the parameter is a string
func (c ConfigurationDirective) ToSeclang() string {
	result := ""
	if c.Metadata != nil {
		result += c.Metadata.ToSeclang()
	}
	result += string(c.Name) + " " + c.Parameter
	return result + "\n"
}
