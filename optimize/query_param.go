package optimize

import (
	"fmt"
	"github.com/moonliightz/go-billwerk/pkg/request"
)

// QueryParam represents a query parameter that can be set on a request.
type QueryParam string

const (
	Range                     QueryParam = "range"
	From                      QueryParam = "from"
	To                        QueryParam = "to"
	Interval                  QueryParam = "interval"
	Size                      QueryParam = "size"
	NextPageToken             QueryParam = "next_page_token"
	Handle                    QueryParam = "handle"
	HandlePrefix              QueryParam = "handle_prefix"
	Handles                   QueryParam = "handles"
	State                     QueryParam = "state"
	ScheduleType              QueryParam = "schedule_type"
	PartialPeriodHandling     QueryParam = "partial_period_handling"
	SetupFeeHandling          QueryParam = "setup_fee_handling"
	FixedLifeTimeUnit         QueryParam = "fixed_life_time_unit"
	TrialIntervalUnit         QueryParam = "trial_interval_unit"
	DunningPlanHandle         QueryParam = "dunning_plan_handle"
	Name                      QueryParam = "name"
	Description               QueryParam = "description"
	SetupFeeText              QueryParam = "setup_fee_text"
	Amount                    QueryParam = "amount"
	Quantity                  QueryParam = "quantity"
	FixedCount                QueryParam = "fixed_count"
	FixedLifeTimeLength       QueryParam = "fixed_life_time_length"
	TrialIntervalLength       QueryParam = "trial_interval_length"
	IntervalLength            QueryParam = "interval_length"
	ScheduleFixedDay          QueryParam = "schedule_fixed_day"
	RenewalReminderEmailDays  QueryParam = "renewal_reminder_email_days"
	TrialReminderEmailDays    QueryParam = "trial_reminder_email_days"
	BaseMonth                 QueryParam = "base_month"
	NoticePeriods             QueryParam = "notice_periods"
	MinimumProratedAmount     QueryParam = "minimum_prorated_amount"
	FixationPeriods           QueryParam = "fixation_periods"
	SetupFee                  QueryParam = "setup_fee"
	AmountInclVAT             QueryParam = "amount_incl_vat"
	NoticePeriodsAfterCurrent QueryParam = "notice_periods_after_current"
	FixationPeriodsFull       QueryParam = "fixation_periods_full"
	IncludeZeroAmount         QueryParam = "include_zero_amount"
	PartialProrationDays      QueryParam = "partial_proration_days"
	FixedTrialDays            QueryParam = "fixed_trial_days"
	Currency                  QueryParam = "currency"
	TaxRateForCountry         QueryParam = "tax_rate_for_country"
)

// QueryParamFunc is a function that sets query parameters on the request builder.
type QueryParamFunc func(requestBuilder request.Builder)

// WithQueryParam adds a single value for a given query parameter.
//
// Example:
//
//	WithQueryParam(Amount, 100) // Results in ?amount=100
func WithQueryParam(param QueryParam, value interface{}) QueryParamFunc {
	return func(requestBuilder request.Builder) {
		requestBuilder.WithParam(string(param), fmt.Sprintf("%v", value))
	}
}

// WithQueryParams adds multiple values for a given query parameter.
//
// Example:
//
//	WithQueryParams(Amount, 100, 200, 300) // Results in ?amount=100&amount=200&amount=300
func WithQueryParams(param QueryParam, values ...interface{}) QueryParamFunc {
	return func(requestBuilder request.Builder) {
		for _, val := range values {
			requestBuilder.AddParam(string(param), fmt.Sprintf("%v", val))
		}
	}
}
