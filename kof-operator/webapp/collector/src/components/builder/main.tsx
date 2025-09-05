import { K8sObject } from "@/models/k8sObject";
import { K8sObjectSet } from "@/models/k8sObjectSet";
import { useMultiClusterServiceProvider } from "@/providers/MultiClusterServicesProvider";
import { ProviderStore } from "@/providers/ProviderAbstract";
import { LucideProps, ThumbsUp } from "lucide-react";
import { ForwardRefExoticComponent, JSX, RefAttributes } from "react";
import { StoreApi, UseBoundStore } from "zustand";
import {
  MultiClusterService,
  MultiClusterServiceSet,
} from "../pages/multi_cluster_services_page/models";
import { CustomizedTableHeadProps } from "../pages/collectorPage/components/collector-list/CollectorTableHead";
import HealthBadge from "../shared/HealthBadge";
import { formatTime } from "@/utils/formatter";
import StatusTab from "../shared/tabs/StatusTab";
import MetadataTab from "../shared/tabs/MetadataTab";
import RawJsonTab from "../shared/tabs/RawJsonTab";

export interface Details<Item> {
  name: string;
  triggerValue: string;
  content: (item: Item) => JSX.Element;
}

export interface DashboardData<Items, Item> {
  name: string;
  id: string;
  store: UseBoundStore<StoreApi<ProviderStore<Items, Item>>>;
  icon: ForwardRefExoticComponent<
    Omit<LucideProps, "ref"> & RefAttributes<SVGSVGElement>
  >;
  tableHeads?: CustomizedTableHeadProps[];
  tableValue?: ((item: Item) => JSX.Element)[];
  details: Details<Item>[];
}

type AnyK8sObject = K8sObject<unknown, unknown, unknown, unknown>;

export const Dashboards: DashboardData<
  K8sObjectSet<AnyK8sObject>,
  AnyK8sObject
>[] = [];

const CreateDashboard1 = (): DashboardData<
  MultiClusterServiceSet,
  MultiClusterService
> => {
  return {
    name: "test",
    id: "test-v1",
    store: useMultiClusterServiceProvider,
    icon: ThumbsUp,
    tableHeads: [
      {
        text: "Name",
        width: 200,
      },
      {
        text: "Status",
        width: 110,
      },
      {
        text: "hello",
        width: 228,
      },
      {
        text: "Age",
        width: 120,
      },
    ],
    tableValue: [
      (item: MultiClusterService) => <>{item.name}</>,
      (item: MultiClusterService) => <HealthBadge isHealthy={item.isHealthy} />,
      (item: MultiClusterService) => <>{item.name}</>,
      (item: MultiClusterService) => <>{formatTime(item.ageInSeconds)}</>,
    ],
    details: [
      {
        name: "Status",
        triggerValue: "status",
        content: (item: MultiClusterService) => (
          <StatusTab conditions={item.status.conditions} />
        ),
      },
      {
        name: "Metadata",
        triggerValue: "metadata",
        content: (item: MultiClusterService) => (
          <MetadataTab metadata={item.metadata} />
        ),
      },
      {
        name: "Raw Json",
        triggerValue: "raw_json",
        content: (item: MultiClusterService) => (
          <RawJsonTab depthLevel={4} object={item.raw} />
        ),
      },
    ],
  };
};

Dashboards.push(CreateDashboard1());
