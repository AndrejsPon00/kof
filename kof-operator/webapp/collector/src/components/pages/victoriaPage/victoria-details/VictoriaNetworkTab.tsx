import { TabsContent } from "@/components/generated/ui/tabs";
import { MetricCardRow, MetricsCard } from "@/components/shared/MetricsCard";
import { VICTORIA_METRICS } from "@/constants/metrics.constants";
import { useVictoriaMetricsState } from "@/providers/victoria_metrics/VictoriaMetricsProvider";
import { bytesToUnits, formatNumber } from "@/utils/formatter";
import { AlertTriangle, Database, Globe, Network, Server } from "lucide-react";
import { JSX } from "react";
import { Pod } from "../../collectorPage/models";

const VictoriaNetworkTab = (): JSX.Element => {
  return (
    <TabsContent value="network" className="flex flex-col gap-5">
      <div className="grid gap-6 lg:grid-cols-3">
        <TCPConnectionDetailsCard />
        <HTTPPerformanceDetailsCard />
        <NetworkErrorsCard />
      </div>
      <div className="grid gap-6 lg:grid-cols-2">
        <BackendConnectionsCard />
        <VictoriaLogsNetworkCard />
      </div>
    </TabsContent>
  );
};

export default VictoriaNetworkTab;

const TCPConnectionDetailsCard = (): JSX.Element => {
  const rows: MetricCardRow[] = [
    {
      title: "Data Read",
      metricName: VICTORIA_METRICS.VM_TCPLISTENER_READ_BYTES_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => bytesToUnits(value),
    },
    {
      title: "Data Written",
      metricName: VICTORIA_METRICS.VM_TCPLISTENER_WRITTEN_BYTES_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => bytesToUnits(value),
    },
    {
      title: "Total Accepts",
      metricName: VICTORIA_METRICS.VM_TCPLISTENER_ACCEPTS_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "Active Connections",
      metricName: VICTORIA_METRICS.VM_TCPLISTENER_CONNS,
      metricFormat: (value: number) => formatNumber(value),
    },
  ];

  return (
    <MetricsCard
      rows={rows}
      icon={Network}
      state={useVictoriaMetricsState}
      title={"TCP Connection Details"}
    />
  );
};

const HTTPPerformanceDetailsCard = (): JSX.Element => {
  const rows: MetricCardRow[] = [
    {
      title: "Total HTTP Requests",
      metricName: VICTORIA_METRICS.VM_HTTP_REQUESTS_ALL_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "Request Errors",
      metricName: VICTORIA_METRICS.VM_HTTP_REQUEST_ERRORS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "Timeout Closed Conns",
      metricName: VICTORIA_METRICS.VM_HTTP_CONN_TIMEOUT_CLOSED_CONNS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "Avg Response Time",
      metricFormat: (value: number) => `${value.toFixed(2)}ms`,
      metricFetchFn: (pod: Pod): number => {
        const requestDurationSec = pod.getMetric(
          VICTORIA_METRICS.VM_HTTP_REQUEST_DURATION_SECONDS_SUM
        );
        const requestDurationCount = pod.getMetric(
          VICTORIA_METRICS.VM_HTTP_REQUEST_DURATION_SECONDS_COUNT
        );
        return (requestDurationSec / requestDurationCount) * 1000;
      },
    },
  ];

  return (
    <MetricsCard
      rows={rows}
      icon={Globe}
      state={useVictoriaMetricsState}
      title={"HTTP Performance Details"}
    />
  );
};

const NetworkErrorsCard = (): JSX.Element => {
  const rows: MetricCardRow[] = [
    {
      title: "VL HTTP Errors",
      metricName: VICTORIA_METRICS.VL_HTTP_ERRORS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "VL UDP Errors",
      metricName: VICTORIA_METRICS.VL_UDP_ERRORS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "TCP Listener Errors",
      metricName: VICTORIA_METRICS.VM_TCPLISTENER_ERRORS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "TCP Read Timeouts",
      metricName: VICTORIA_METRICS.VM_TCPLISTENER_READ_TIMEOUTS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
  ];

  return (
    <MetricsCard
      rows={rows}
      icon={AlertTriangle}
      state={useVictoriaMetricsState}
      title={"Network Errors & Timeouts"}
    />
  );
};

const BackendConnectionsCard = (): JSX.Element => {
  const rows: MetricCardRow[] = [
    {
      title: "Backend Bytes Read",
      metricName: VICTORIA_METRICS.VLSELECT_BACKEND_CONN_BYTES_READ_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => bytesToUnits(value),
    },
    {
      title: "Backend Bytes Written",
      metricName: VICTORIA_METRICS.VLSELECT_BACKEND_CONN_BYTES_WRITTEN_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => bytesToUnits(value),
    },
    {
      title: "Backend Read Errors",
      metricName: VICTORIA_METRICS.VLSELECT_BACKEND_CONN_READ_ERRORS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "Active Backend Conns",
      metricName: VICTORIA_METRICS.VLSELECT_BACKEND_CONNS,
      metricFormat: (value: number) => formatNumber(value),
    },
  ];

  return (
    <MetricsCard
      rows={rows}
      icon={Server}
      state={useVictoriaMetricsState}
      title={"VLSelect Backend Connections"}
    />
  );
};

const VictoriaLogsNetworkCard = (): JSX.Element => {
  const rows: MetricCardRow[] = [
    {
      title: "VL HTTP Requests",
      metricName: VICTORIA_METRICS.VL_HTTP_REQUESTS_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "VL UDP Requests",
      metricName: VICTORIA_METRICS.VL_UDP_REQESTS_TOTAL,
      enableTrendSystem: true,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "VL All Errors",
      metricName: VICTORIA_METRICS.VL_ERRORS_TOTAL,
      enableTrendSystem: true,
      isPositiveTrend: false,
      metricFormat: (value: number) => formatNumber(value),
    },
    {
      title: "VL Request Duration",
      metricName: VICTORIA_METRICS.VL_HTTP_REQUEST_DURATION_SECONDS_SUM,
      metricFormat: (value: number) => `${value.toFixed(3)}s`,
    },
  ];

  return (
    <MetricsCard
      rows={rows}
      icon={Database}
      state={useVictoriaMetricsState}
      title={"VictoriaLogs Network"}
    />
  );
};
