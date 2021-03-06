package ecode

var (
	GameGetUniqueIdFailed         = New(14000) //获取unique id失败
	GameConnectServiceFail        = New(14001)
	GameConnectDbFail             = New(14002)
	GameDbTimeout                 = New(14003)
	GameCreateFail                = New(14004)
	GameIdIsEmpty                 = New(14005)
	GameRecordNotFound            = New(14006)
	GameOrderValueIsEmpty         = New(14007)
	GameUpdateFail                = New(14008)
	GameUpdateStatusFail          = New(14009)
	GameGetDetailFail             = New(14010)
	GameCreateEquipmentFail       = New(14011)
	GameUpdateEquipmentFail       = New(14012)
	GameUpdateEquipmentStatusFail = New(14013)
	GameGetEquipmentDetailFail    = New(14014)
	GameCreateRoleFail            = New(14015)
	GameUpdateRoleFail            = New(14016)
	GameUpdateRoleStatusFail      = New(14017)
	GameGetRoleDetailFail         = New(14018)
	GameCreateSkillFail           = New(14019)
	GameUpdateSkillFail           = New(14020)
	GameUpdateSkillStatusFail     = New(14021)
	GameGetSkillDetailFail        = New(14022)
)
