import { Clock } from "lucide-react";
import React from "react";
import IconText from "~/components/texts/IconText";

export default function DurationText({
  startDate,
  endDate,
  long,
}: {
  startDate: Date;
  endDate?: Date;
  long?: boolean;
}) {
  const [now, setNow] = React.useState(new Date());

  React.useEffect(() => {
    const intervalId = setInterval(() => {
      setNow(new Date());
    }, 1000);

    return () => clearInterval(intervalId);
  }, []);

  const [start, setStart] = React.useState(startDate);
  const [end, setEnd] = React.useState(endDate ?? now);

  React.useEffect(() => {
    setStart(startDate);
    setEnd(endDate ?? now);
  }, [startDate, endDate, now]);

  const [duration, setDuration] = React.useState(
    () => end.getTime() - start.getTime()
  );

  React.useEffect(() => {
    setDuration(end.getTime() - start.getTime());
  }, [start, end]);

  const [days, hours, minutes, seconds] = React.useMemo(() => {
    const totalSeconds = Math.floor(duration / 1000);
    const secs = totalSeconds % 60;
    const mins = Math.floor(totalSeconds / 60) % 60;
    const hrs = Math.floor(totalSeconds / 3600) % 24;
    const d = Math.floor(totalSeconds / 86400);
    return [d, hrs, mins, secs];
  }, [duration]);

  const pluralize = (value: number, singular: string, plural: string) =>
    value === 1
      ? `${value} ${singular}`
      : value !== 0
      ? `${value} ${plural}`
      : "";

  const daysPart = pluralize(days, "day", "days");
  const hoursPart = pluralize(hours, "hour", "hours");
  const minutesPart = pluralize(minutes, "minute", "minutes");
  const secondsPart = pluralize(seconds, "second", "seconds");

  const longParts = [];
  if (daysPart) longParts.push(daysPart);
  if (hoursPart) longParts.push(hoursPart);
  if (minutesPart) longParts.push(minutesPart);
  if (secondsPart) longParts.push(secondsPart);
  const durationTextLong = longParts.join(", ");

  const daysShort = days !== 0 ? days + "d" : "";
  const hoursShort = hours !== 0 ? hours + "h" : "";
  const minutesShort = minutes !== 0 ? minutes + "m" : "";
  const secondsShort = seconds !== 0 ? seconds + "s" : "";

  const shortParts = [];
  if (daysShort) shortParts.push(daysShort);
  if (hoursShort) shortParts.push(hoursShort);
  if (minutesShort) shortParts.push(minutesShort);
  if (secondsShort) shortParts.push(secondsShort);
  const durationTextShort = shortParts.join(" ");

  const durationText = long ? durationTextLong : durationTextShort;

  return <IconText icon={Clock} text={durationText} />;
}
