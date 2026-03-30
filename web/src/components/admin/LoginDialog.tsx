import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/authStore";

interface LoginDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function LoginDialog({ open, onOpenChange }: LoginDialogProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { login } = useAuthStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login(username, password);
      onOpenChange(false);
      setUsername("");
      setPassword("");
    } catch {
      setError("Invalid username or password");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="border-[var(--border-color)] bg-[var(--card-bg)] w-[calc(100vw-2rem)] sm:w-auto sm:max-w-sm rounded-xl">
        <DialogHeader>
          <DialogTitle className="text-[var(--text-primary)]">Admin Login</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label className="text-[var(--text-secondary)]">Username</Label>
            <Input
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="border-[var(--border-color)] bg-[var(--card-bg-secondary)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)]"
              placeholder="Enter username"
              autoFocus
            />
          </div>
          <div className="space-y-2">
            <Label className="text-[var(--text-secondary)]">Password</Label>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="border-[var(--border-color)] bg-[var(--card-bg-secondary)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)]"
              placeholder="Enter password"
            />
          </div>
          {error && <p className="text-sm text-[#E57373]">{error}</p>}
          <Button
            type="submit"
            className="w-full"
            disabled={loading || !username || !password}
            style={{ backgroundColor: "var(--theme-primary)" }}
          >
            {loading ? "Signing in..." : "Sign In"}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
