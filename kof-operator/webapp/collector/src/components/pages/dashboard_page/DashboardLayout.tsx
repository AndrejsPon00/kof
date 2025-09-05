import { JSX } from "react";
import { Outlet } from "react-router-dom";
import { Separator } from "@/components/generated/ui/separator";
import Loader from "@/components/shared/Loader";
import FetchStatus from "@/components/shared/FetchStatus";
import { DashboardData } from "@/components/builder/main";
import { K8sObjectSet } from "@/models/k8sObjectSet";
import { K8sObject } from "@/models/k8sObject";

const DashboardLayout = <
  Items extends K8sObjectSet<Item> | unknown[],
  Item extends K8sObject<unknown, unknown, unknown, unknown>
>({
  name,
  store,
}: DashboardData<Items, Item>): JSX.Element => {
  const { items, isLoading, error, fetch } = store();

  return (
    <div className="flex flex-col w-full h-full p-5 space-y-8">
      <h1 className="font-bold text-3xl">{name}</h1>
      <Separator />
      {isLoading ? (
        <Loader />
      ) : error ? (
        <FetchStatus onReload={fetch}>
          Failed to fetch {name}. Click "Reload" button to try again.
        </FetchStatus>
      ) : !items || !items.length ? (
        <FetchStatus onReload={fetch}>No {name} found</FetchStatus>
      ) : (
        <Outlet />
      )}
    </div>
  );
};

export default DashboardLayout;
