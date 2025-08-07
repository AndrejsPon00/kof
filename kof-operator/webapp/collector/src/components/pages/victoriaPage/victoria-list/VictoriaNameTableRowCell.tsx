import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/generated/ui/hover-card";
import { TableCell } from "@/components/generated/ui/table";
import { JSX } from "react";

type Name =
  | "vlselect"
  | "vlinsert"
  | "vlstorage"
  | "vmselect"
  | "vminsert"
  | "vmstorage";

const VICTORIA_TYPES: Name[] = [
  "vlselect",
  "vlinsert",
  "vlstorage",
  "vmselect",
  "vminsert",
  "vmstorage",
];

const NAME_TYPE: { [name in Name]: string } = {
  vlselect: "VictoriaLogs Select",
  vlinsert: "VictoriaLogs Insert",
  vlstorage: "VictoriaLogs Storage",
  vmselect: "VictoriaMetrics Select",
  vminsert: "VictoriaMetrics Insert",
  vmstorage: "VictoriaMetrics Storage",
};

const VictoriaNameTableRowCell = ({ name }: { name: string }): JSX.Element => {
  const victoriaType: Name | undefined = VICTORIA_TYPES.find((type) => {
    return name.includes(type);
  });

  return (
    <TableCell>
      <HoverCard>
        <HoverCardTrigger className="flex flex-col">
          <span className="truncate">{name}</span>
          <span className="text-sm text-muted-foreground">
            {victoriaType ? NAME_TYPE[victoriaType] : ""}
          </span>
        </HoverCardTrigger>
        <HoverCardContent className="w-fit">{name}</HoverCardContent>
      </HoverCard>
    </TableCell>
  );
};

export default VictoriaNameTableRowCell;
