import { formatTime } from "@/utils/formatter";
import { ItemObject, TableColumn } from "./DashboardTypes";
import HealthBadge from "@/components/shared/HealthBadge";

export const TableColNamespace = <Item extends ItemObject>(): TableColumn<Item> => {
  return {
    head: { text: "Namespace", width: 110 },
    valueFn: (item: Item) => <>{item.namespace}</>,
  };
};

export const TableColName = <Item extends ItemObject>(): TableColumn<Item> => {
  return {
    head: { text: "Name", width: 200 },
    valueFn: (item: Item) => <>{item.name}</>,
  };
};

export const TableColAge = <Item extends ItemObject>(): TableColumn<Item> => {
  return {
    head: { text: "Age", width: 120 },
    valueFn: (item: Item) => <>{formatTime(item.ageInSeconds)}</>,
  };
};

export const TableColStatus = <Item extends ItemObject>(): TableColumn<Item> => {
  return {
    head: { text: "Status", width: 110 },
    valueFn: (item: Item) => <HealthBadge isHealthy={item.isHealthy} />,
  };
};
