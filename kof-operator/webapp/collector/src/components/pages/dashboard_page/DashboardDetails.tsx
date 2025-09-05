import { Button } from "@/components/generated/ui/button";
import { Tabs, TabsList, TabsTrigger } from "@/components/generated/ui/tabs";
import { Layers2, MoveLeft } from "lucide-react";
import { JSX, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { DashboardData } from "@/components/builder/main";
import { K8sObjectSet } from "@/models/k8sObjectSet";
import { K8sObject } from "@/models/k8sObject";
import MetadataTab from "@/components/shared/tabs/MetadataTab";
import RawJsonTab from "@/components/shared/tabs/RawJsonTab";
import StatusTab from "@/components/shared/tabs/StatusTab";
import DetailsHeader from "@/components/shared/DetailsHeader";

const DashboardDetails = <
  Items extends K8sObjectSet<Item> | unknown[],
  Item extends K8sObject<unknown, unknown, unknown, unknown>
>({
  name,
  store,
  details,
}: DashboardData<Items, Item>): JSX.Element => {
  const { isLoading, items, selectedItem, selectItem } = store();

  const navigate = useNavigate();
  const { objName } = useParams();

  useEffect(() => {
    if (!isLoading && items && objName) {
      selectItem(objName);
    }
  }, [objName, items, isLoading, selectItem]);

  if (!selectedItem) {
    return (
      <div className="flex flex-col w-full h-full p-5 space-y-8">
        <div className="flex flex-col w-full h-full justify-center items-center space-y-4">
          <span className="font-bold text-2xl">{name} not found</span>
          <Button
            variant="outline"
            className="cursor-pointer"
            onClick={() => {
              navigate(-1);
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
    <div className="flex flex-col gap-5">
      <DetailsHeader
        icon={Layers2}
        title={selectedItem.name}
        isHealthy={selectedItem.isHealthy}
      />
      <Tabs defaultValue="status" className="space-y-6">
        <TabsList className="flex w-full">
          {details?.map((d) => (
            <TabsTrigger value={d.triggerValue}>{d.name}</TabsTrigger>
          ))}
          {/* <TabsTrigger value="status">Status</TabsTrigger>
          <TabsTrigger value="metadata">Metadata</TabsTrigger>
          <TabsTrigger value="raw_json">Raw Json</TabsTrigger> */}
        </TabsList>
        {details?.map((d) => d.content(selectedItem))}
        {/* <StatusTab conditions={selectedItem.status.conditions} />
        <MetadataTab metadata={selectedItem.metadata} />
        <RawJsonTab depthLevel={4} object={selectedItem.raw} /> */}
      </Tabs>
    </div>
  );
};

export default DashboardDetails;
