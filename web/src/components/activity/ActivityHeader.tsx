import type { ActivitySettings } from "@/types";

interface ActivityHeaderProps {
  settings: ActivitySettings;
}

export function ActivityHeader({ settings }: ActivityHeaderProps) {
  return (
    <div className="rounded-xl border border-[var(--border-color)] glass-card p-3 md:p-4">
      <div className="flex items-center gap-3 md:gap-4">
        {settings.logo_path && (
          <img
            src={settings.logo_path}
            alt="Logo"
            className="max-h-10 md:max-h-14 max-w-24 md:max-w-32 shrink-0 rounded-lg object-contain"
          />
        )}
        <h1 className="text-lg md:text-2xl font-bold tracking-tight text-[var(--text-primary)] truncate">
          {settings.activity_name}
        </h1>
      </div>
    </div>
  );
}
