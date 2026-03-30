import { useState, useEffect } from "react";
import { Cloud, Sun, CloudRain, CloudSnow, CloudLightning, CloudFog } from "lucide-react";

interface WeatherData {
  temp: number;
  condition: string;
  icon: string;
}

export function DateTimeWeather() {
  const [time, setTime] = useState(new Date());
  const [weather, setWeather] = useState<WeatherData | null>(null);

  // Update time every second
  useEffect(() => {
    const timer = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  // Fetch weather on mount (using a free API with IP-based location)
  useEffect(() => {
    fetch("https://wttr.in/?format=j1")
      .then((res) => res.json())
      .then((data) => {
        const current = data.current_condition?.[0];
        if (current) {
          setWeather({
            temp: parseInt(current.temp_C),
            condition: current.weatherDesc?.[0]?.value || "Unknown",
            icon: current.weatherCode,
          });
        }
      })
      .catch(() => {
        // Silently fail - weather is not critical
      });
  }, []);

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    });
  };

  const formatDate = (date: Date) => {
    return date.toLocaleDateString("en-US", {
      weekday: "short",
      month: "short",
      day: "numeric",
    });
  };

  const getWeatherIcon = (code: string) => {
    const codeNum = parseInt(code);
    // Weather codes from wttr.in
    if (codeNum >= 200 && codeNum < 300) return <CloudLightning className="h-5 w-5" />;
    if (codeNum >= 300 && codeNum < 600) return <CloudRain className="h-5 w-5" />;
    if (codeNum >= 600 && codeNum < 700) return <CloudSnow className="h-5 w-5" />;
    if (codeNum >= 700 && codeNum < 800) return <CloudFog className="h-5 w-5" />;
    if (codeNum === 800) return <Sun className="h-5 w-5" />;
    return <Cloud className="h-5 w-5" />;
  };

  return (
    <div className="rounded-xl border border-[var(--border-color)] glass-card p-3 md:p-4">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-xl md:text-2xl font-bold tabular-nums text-[var(--text-primary)]">
            {formatTime(time)}
          </div>
          <div className="text-xs md:text-sm text-[var(--text-secondary)]">
            {formatDate(time)}
          </div>
        </div>
        {weather && (
          <div className="flex items-center gap-2 text-[var(--text-secondary)]">
            {getWeatherIcon(weather.icon)}
            <span className="text-lg font-medium text-[var(--text-primary)]">
              {weather.temp}°C
            </span>
          </div>
        )}
      </div>
    </div>
  );
}
