import { Moon, Sun } from "lucide-react";
import { Outlet } from "react-router-dom";
import { Button } from "~/components/ui/button";
import { cn } from "~/lib/utils";
import { useTheme } from "~/providers/ThemeProvider";

const TopBar: React.FC = () => {
  const { theme, setTheme } = useTheme();
  const toggleTheme = () => {
    setTheme(theme === "light" ? "dark" : "light");
  };
  const ThemeIcon = theme === "light" ? Moon : Sun;

  return (
    <header
      className={cn(
        "sticky top-0 z-30 no-select",
        "flex items-center justify-between h-0 px-6 shadow-sm",
        "translate-y-8"
      )}
    >
      <div></div>
      <Button
        variant="ghost"
        size="icon"
        onClick={toggleTheme}
        aria-label="Toggle theme"
        className="bg-transparent"
      >
        <ThemeIcon className="h-5 w-5" />
      </Button>
    </header>
  );
};

export default function AdminLayout() {
  return (
    <>
      <div className="background-gradient" />
      <TopBar />
      <div className="flex min-h-screen flex-col items-center justify-center background">
        <Outlet />
      </div>
    </>
  );
}
