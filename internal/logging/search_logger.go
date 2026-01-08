package logging

import "log/slog"

type ExecutionContext struct {
	RunID   string
	QueryID string
}

func LogOriginalQuery(
	logger *slog.Logger,
	ctx ExecutionContext,
	query string,
) {
	logger.Debug("original query",
		slog.String("run_id", ctx.RunID),
		slog.String("query_id", ctx.QueryID),
		slog.String("query", query),
	)
}

func LogEnhancedQuery(
	logger *slog.Logger,
	ctx ExecutionContext,
	data EnhancedQueryLog,
) {
	logger.Debug("enhanced query",
		slog.String("run_id", ctx.RunID),
		slog.String("query_id", ctx.QueryID),
		slog.String("enhancement", data.EnhancementType),
		slog.String("original_query", data.OriginalQuery),
		slog.String("enhanced_query", data.EnhancedQuery),
	)
}

func LogRRFResults(
	logger *slog.Logger,
	ctx ExecutionContext,
	candidates []RRFCandidateLog,
) {
	logger.Debug("rrf results",
		slog.String("run_id", ctx.RunID),
		slog.Int("candidate_count", len(candidates)),
		slog.Any("candidates", candidates),
	)
}

func LogFinalResults(
	logger *slog.Logger,
	ctx ExecutionContext,
	results []FinalResultLog,
) {
	logger.Debug("final ranked results",
		slog.String("run_id", ctx.RunID),
		slog.Int("result_count", len(results)),
		slog.Any("results", results),
	)
}
