import { JSX } from "react";
import { Pod } from "../../models";
import { useCollectorMetricsState } from "@/providers/collectors_metrics/CollectorsMetricsProvider";
import { useTimePeriod } from "@/providers/collectors_metrics/TimePeriodState";
import { getMetricTrendData } from "@/utils/metrics";
import { formatNumber } from "@/utils/formatter";
import { TableCell } from "@/components/generated/ui/table";
import { TrendingUp } from "lucide-react";

interface MetricTrendCellProps {
  metric: string;
  pod: Pod;
}

const MetricTrendCell = ({
  pod,
  metric,
}: MetricTrendCellProps): JSX.Element => {
  const { metricsHistory } = useCollectorMetricsState();
  const { timePeriod } = useTimePeriod();

  const { metricValue, metricTrend } = getMetricTrendData(
    metric,
    metricsHistory,
    pod,
    timePeriod
  );

  const formattedMetricValue = formatNumber(metricValue);

  const trendMessageColor = metricTrend.isTrending
    ? "text-green-600"
    : "text-red-600";

  return (
    <TableCell>
      <div className="flex flex-col">
        {metricTrend && (
          <div
            className={`flex gap-2 items-center font-semibold ${trendMessageColor}`}
          >
            {metricTrend.isTrending && <TrendingUp className="w-4 h-4" />}
            {metricTrend.message}
          </div>
        )}
        <span className={`text-xs`}>Total: {formattedMetricValue}</span>
      </div>
    </TableCell>
  );
};

export default MetricTrendCell;
