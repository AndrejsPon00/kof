import { useVictoriaMetricsState } from "@/providers/victoria_metrics/VictoriaMetricsProvider";
import { JSX, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import VictoriaPageHeader from "../VictoriaPageHeader";
import { Separator } from "@/components/generated/ui/separator";
import { Loader, MoveLeft } from "lucide-react";
import { Button } from "@/components/generated/ui/button";
import UnhealthyAlert from "../../collectorPage/components/collector-details/UnhealthyAlert";
import { Tabs, TabsList, TabsTrigger } from "@/components/generated/ui/tabs";
import VictoriaOverviewTab from "./VictoriaOverviewTab";
import VictoriaSystemTab from "./VictoriaSystemTab";
import VictoriaGoRuntimeTab from "./VictoriaGoRuntimeTab";
import VictoriaNetworkTab from "./VictoriaNetworkTab";

const VictoriaDetailsPage = (): JSX.Element => {
  const { setSelectedPod, setSelectedCluster, selectedPod, isLoading, data } =
    useVictoriaMetricsState();

  const navigate = useNavigate();

  const { cluster, pod } = useParams();
  useEffect(() => {
    if (!isLoading && cluster && pod) {
      setSelectedCluster(cluster);
      setSelectedPod(pod);
    }
  }, [cluster, pod, setSelectedCluster, setSelectedPod, isLoading]);

  if (isLoading && !data) {
    return (
      <div className="flex flex-col w-full h-full p-5 space-y-8">
        <VictoriaPageHeader />
        <Separator />
        <div className="flex w-full h-full justify-center items-center">
          <Loader className="animate-spin w-8 h-8"></Loader>
        </div>
      </div>
    );
  }

  if (!selectedPod) {
    return (
      <div className="flex flex-col w-full h-full p-5 space-y-8">
        <VictoriaPageHeader />
        <Separator />
        <div className="flex flex-col w-full h-full justify-center items-center space-y-4">
          <span className="font-bold text-2xl">Collector not found</span>
          <Button
            variant="outline"
            className="cursor-pointer"
            onClick={() => {
              navigate("/collectors");
            }}
          >
            <MoveLeft />
            <span>Back to Table</span>
          </Button>
        </div>
      </div>
    );
  }
  return (
    <div className="flex flex-col w-full h-full p-5 space-y-8">
      <VictoriaPageHeader />
      <Separator />
      {/* <CollectorContentHeader /> */}

      <UnhealthyAlert />
      <Tabs defaultValue="overview" className="space-y-6">
        <TabsList className="flex w-full cursor-pointer">
          <TabsTrigger className="cursor-pointer" value="overview">
            Overview
          </TabsTrigger>
          <TabsTrigger className="cursor-pointer" value="system">
            System
          </TabsTrigger>
          <TabsTrigger className="cursor-pointer" value="network">
            Network
          </TabsTrigger>
          {/* <TabsTrigger className="cursor-pointer" value="victoria_metrics">
            Victoria Metrics
          </TabsTrigger> */}
          <TabsTrigger className="cursor-pointer" value="go_runtime">
            Go Runtime
          </TabsTrigger>
        </TabsList>
        <VictoriaOverviewTab />
        <VictoriaSystemTab />
        <VictoriaNetworkTab />
        <VictoriaGoRuntimeTab />
      </Tabs>
    </div>
  );
};

export default VictoriaDetailsPage;
