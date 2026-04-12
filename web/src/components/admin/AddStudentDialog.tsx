import { useState, useEffect, useRef, useCallback } from "react";
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
import { searchAutocomplete, type AutocompleteEntry } from "@/api/autocomplete";
import type { House } from "@/types";
import api from "@/api/client";

interface AddStudentDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function AddStudentDialog({ open, onOpenChange, onSuccess }: AddStudentDialogProps) {
  const [name, setName] = useState("");
  const [studentId, setStudentId] = useState("");
  const [houseId, setHouseId] = useState("");
  const [lapCount, setLapCount] = useState("");
  const [houses, setHouses] = useState<House[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [suggestions, setSuggestions] = useState<AutocompleteEntry[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [activeField, setActiveField] = useState<"name" | "studentId" | null>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(null);

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
    } else {
      setSuggestions([]);
      setShowSuggestions(false);
    }
  }, [open]);

  const doSearch = useCallback((query: string) => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (!query.trim()) {
      setSuggestions([]);
      setShowSuggestions(false);
      return;
    }
    debounceRef.current = setTimeout(async () => {
      const results = await searchAutocomplete(query.trim()).catch(() => []);
      setSuggestions(results);
      setShowSuggestions(results.length > 0);
    }, 300);
  }, []);

  const handleSelectSuggestion = (entry: AutocompleteEntry) => {
    setName(entry.system_name);
    setStudentId(entry.id);
    setSuggestions([]);
    setShowSuggestions(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !houseId) return;

    setError("");
    setLoading(true);
    try {
      await createStudent(name.trim(), parseInt(houseId), parseInt(lapCount) || 0, studentId.trim() || undefined);
      onOpenChange(false);
      setName("");
      setStudentId("");
      setHouseId("");
      setLapCount("");
      onSuccess?.();
    } catch {
      setError("Failed to add student, please try again");
    } finally {
      setLoading(false);
    }
  };

  const renderSuggestions = () => {
    if (!showSuggestions || suggestions.length === 0) return null;
    return (
      <div
        ref={suggestionsRef}
        className="absolute z-50 top-full left-0 right-0 mt-1 max-h-48 overflow-y-auto rounded-md border border-[var(--border-color)] bg-[var(--card-bg)] shadow-lg"
      >
        {suggestions.map((entry) => (
          <button
            key={entry.id}
            type="button"
            className="w-full px-3 py-2 text-left hover:bg-[var(--card-bg-secondary)] transition-colors flex items-center justify-between gap-2"
            onMouseDown={(e) => {
              e.preventDefault();
              handleSelectSuggestion(entry);
            }}
          >
            <span className="text-sm text-[var(--text-primary)] truncate">{entry.system_name}</span>
            <span className="text-xs text-[var(--text-muted)] shrink-0">{entry.id}</span>
          </button>
        ))}
      </div>
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="border-[var(--border-color)] bg-[var(--card-bg)] w-[calc(100vw-2rem)] sm:w-full sm:max-w-md rounded-xl">
        <DialogHeader>
          <DialogTitle className="text-[var(--text-primary)]">Add Student</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2 relative">
            <Label className="text-[var(--text-secondary)]">Name</Label>
            <Input
              value={name}
              onChange={(e) => {
                setName(e.target.value);
                setActiveField("name");
                doSearch(e.target.value);
              }}
              onFocus={() => {
                setActiveField("name");
                if (suggestions.length > 0) setShowSuggestions(true);
              }}
              onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
              className="border-[var(--border-color)] bg-[var(--card-bg-secondary)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)]"
              placeholder="Enter student name"
              autoFocus
              autoComplete="off"
            />
            {activeField === "name" && renderSuggestions()}
          </div>
          <div className="space-y-2 relative">
            <Label className="text-[var(--text-secondary)]">Student ID <span className="text-[var(--text-muted)]">(optional)</span></Label>
            <Input
              value={studentId}
              onChange={(e) => {
                setStudentId(e.target.value);
                setActiveField("studentId");
                doSearch(e.target.value);
              }}
              onFocus={() => {
                setActiveField("studentId");
                if (suggestions.length > 0) setShowSuggestions(true);
              }}
              onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
              className="border-[var(--border-color)] bg-[var(--card-bg-secondary)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)]"
              placeholder="e.g. s23321"
              autoComplete="off"
            />
            {activeField === "studentId" && renderSuggestions()}
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
