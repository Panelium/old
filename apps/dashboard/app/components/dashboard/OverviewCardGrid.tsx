import React from "react";
import type { OverviewCardProps } from "./OverviewCard";
import OverviewCard from "./OverviewCard";

export default function OverviewCardGrid({
  cards,
}: {
  cards: OverviewCardProps[];
}) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      {cards.map((props, idx) => (
        <OverviewCard key={idx} {...props} />
      ))}
    </div>
  );
}
