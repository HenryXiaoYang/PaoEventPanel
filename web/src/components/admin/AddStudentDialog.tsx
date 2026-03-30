import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { createStudent } from "@/api/students";
import type { House } from "@/types";
import api from "@/api/client";

interface AddStudentDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function AddStudentDialog({ open, onOpenChange, onSuccess }: AddStudentDialogProps) {
  const [name, setName] = useState("");
  const [houseId, setHouseId] = useState("");
  const [lapCount, setLapCount] = useState("");
  const [houses, setHouses] = useState<House[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (open) {
      api.get("/rankings/houses").then((res) => {
        const houseList = res.data.map((h: { house_id: number; name: string; code: string; color: string }) => ({
          id: h.house_id,
          name: h.name,
          code: h.code,
          color: h.color,
        }));
        setHouses(houseList);
      });
    }
  }, [open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !houseId) return;

    setError("");
    setLoading(true);
    try {
      await createStudent(name.trim(), parseInt(houseId), parseInt(lapCount) || 0);
      onOpenChange(false);
      setName("");
      setHouseId("");
      setLapCount("");
      onSuccess?.();
    } catch {
      setError("Failed to add student, please try again");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="border-[var(--border-color)] bg-[var(--card-bg)] w-[calc(100vw-2rem)] sm:w-auto sm:max-w-sm rounded-xl">
        <DialogHeader>
          <DialogTitle className="text-[var(--text-primary)]">Add Student</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label className="text-[var(--text-secondary)]">Name</Label>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="border-[var(--border-color)] bg-[var(--card-bg-secondary)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)]"
              placeholder="Enter student name"
              autoFocus
            />
          </div>
          <div className="space-y-2">
            <Label className="text-[var(--text-secondary)]">House</Label>
            <Select value={houseId} onValueChange={(value) => setHouseId(value ?? "")}>
              <SelectTrigger className="border-[var(--border-color)] bg-[var(--card-bg-secondary)] text-[var(--text-primary)]">
                <SelectValue placeholder="Select house" />
              </SelectTrigger>
              <SelectContent className="border-[var(--border-color)] bg-[var(--card-bg)]">
                {houses.map((house) => (
                  <SelectItem
                    key={house.id}
                    value={String(house.id)}
                    className="text-[var(--text-primary)] focus:bg-[var(--card-bg-secondary)] focus:text-[var(--text-primary)]"
                  >
                    <div className="flex items-center gap-2">
                      <div
                        className="h-2.5 w-2.5 rounded-full"
                        style={{ backgroundColor: house.color }}
                      />
                      {house.name}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label className="text-[var(--text-secondary)]">Lap Count</Label>
            <Input
              type="number"
              min="0"
              value={lapCount}
              onChange={(e) => setLapCount(e.target.value)}
              className="border-[var(--border-color)] bg-[var(--card-bg-secondary)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)]"
              placeholder="0"
            />
          </div>
          {error && <p className="text-sm text-[#E57373]">{error}</p>}
          <Button
            type="submit"
            className="w-full"
            disabled={loading || !name.trim() || !houseId}
            style={{ backgroundColor: "var(--theme-primary)" }}
          >
            {loading ? "Adding..." : "Add"}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
