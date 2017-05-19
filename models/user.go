package models

type User struct {
	Id                int64     `json:"id" orm:"column(ID)"`
	Name              string    `json:"name" orm:"column(NAME)"`
	Header            string    `json:"header" orm:"column(HEADER)"`
	Email             string    `json:"email" orm:"column(EMAIL)"`
	Phone             string    `json:"phone" orm:"column(PHONE)"`
	Company           string    `json:"company" orm:"column(COMPANY)"`
	EncryptedPassword string    `json:"encryptedPassword" orm:"column(ENCRYPTED_PASSWORD)"`
	CreatedAt         int64     `json:"createdAt" orm:"column(CREATED_AT)"`
	UpdatedAt         int64     `json:"updatedAt" orm:"column(UPDATED_AT)"`
	Active            bool      `json:"active" orm:"column(ACTIVE)"`
	Role              int       `json:"role" orm:"column(ROLE)"`
	Resource          *Resource `json:"resource" orm:"reverse(one)"`
}

type Resource struct {
	Id                  int64   `json:"id" orm:"column(ID)"`
	AlgorithmResource   string  `json:"algorithmResource" orm:"column(ALGORITHM_RESOURCE);type(text)"`
	CpuTotalResource    float64 `json:"cpuTotalResource" orm:"column(CPU_TOTAL_RESOURCE)"`
	CpuUsageResource    float64 `json:"cpuUsageResource" orm:"column(CPU_USAGE_RESOURCE)"`
	CpuUnit             string  `json:"cpuUnit" orm:"column(CPU_UNIT)"`
	MemoryTotalResource float64 `json:"memoryTotalResource" orm:"column(MEMORY_TOTAL_RESOURCE)"`
	MemoryUsageResource float64 `json:"memoryUsageResource" orm:"column(MEMORY_USAGE_RESOURCE)"`
	MemoryUnit          string  `json:"memoryUnit" orm:"column(MEMORY_UNIT)"`
	User                *User   `json:"user" orm:"column(USER_ID);rel(one)"`
	QuotaNamespace      string  `json:"quotaNamespace" orm:"column(QUOTA_NAMESPACE)"`
	QuotaName           string  `json:"quotaName" orm:"column(QUOTA_NAME)"`
}
