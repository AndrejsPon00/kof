import { JSX } from "react";
import { TabsContent } from "@/components/generated/ui/tabs";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/generated/ui/card";
import { Progress } from "@/components/generated/ui/progress";
import { Separator } from "@/components/generated/ui/separator";
import { METRICS } from "@/constants/metrics.constants";
import { formatNumber } from "@/utils/formatter";
import { useCollectorMetricsState } from "@/providers/collectors_metrics/CollectorsMetricsProvider";
import StatRow from "@/components/shared/StatRow";
import { MetricCardRow, MetricsCard } from "@/components/shared/MetricsCard";
import { Send, TriangleAlert } from "lucide-react";

const CollectorExporterTabContent = (): JSX.Element => {
  return (
    <TabsContent
      value="exporter"
      className="grid gap-6 md:grid-cols-2 lg:grid-cols-3"
    >
      <QueueCard />
      <SentRecordsCard />
      <FailedRecordsCard />
    </TabsContent>
  );
};

const QueueCard = (): JSX.Element => {
  const { selectedPod: pod } = useCollectorMetricsState();

  if (!pod) {
    return <></>;
  }

  const queueCapacity = pod.getMetric(METRICS.OTELCOL_EXPORTER_QUEUE_CAPACITY);
  const queueSize = pod.getMetric(METRICS.OTELCOL_EXPORTER_QUEUE_SIZE);
  const queueUtilization = (queueSize / queueCapacity) * 100;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">Queue Status</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <StatRow text="Capacity" value={queueCapacity} />
        <StatRow text="Current Size" value={queueSize} />
        <StatRow
          text="Utilization"
          value={`${queueUtilization.toFixed(1)}%`}
          valueStyles="text-sm"
          containerStyle="mb-2"
        />
        <Progress value={queueUtilization} />
      </CardContent>
    </Card>
  );
};

const SentRecordsCard = (): JSX.Element => {
  const rows: MetricCardRow[] = [
    {
      title: "Log Records",
      metricName: METRICS.OTELCOL_EXPORTER_SENT_LOG_RECORDS,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "Metric Points",
      metricName: METRICS.OTELCOL_EXPORTER_SENT_METRIC_POINTS,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },

    {
      title: "",
      customRow: (
        <>
          <Separator />
          <div className="text-xs text-muted-foreground">
            Total records successfully exported
          </div>
        </>
      ),
    },
  ];

  return (
    <MetricsCard
      rows={rows}
      icon={Send}
      state={useCollectorMetricsState}
      title={"Sent Records"}
    />
  );
};

const FailedRecordsCard = (): JSX.Element => {
  const rows: MetricCardRow[] = [
    {
      title: "Failed Log Records",
      metricName: METRICS.OTELCOL_EXPORTER_SEND_FAILED_LOG_RECORDS,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "Failed Metric Points",
      metricName: METRICS.OTELCOL_EXPORTER_SEND_FAILED_METRIC_POINTS,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },

    {
      title: "",
      customRow: (
        <>
          <Separator />
          <div className="text-xs text-muted-foreground">
            Records that failed to export
          </div>
        </>
      ),
    },
  ];

  return (
    <MetricsCard
      rows={rows}
      icon={TriangleAlert}
      state={useCollectorMetricsState}
      title={"Failed Records"}
    />
  );
};

export default CollectorExporterTabContent;
