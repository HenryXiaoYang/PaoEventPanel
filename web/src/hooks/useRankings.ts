import { useRef, useMemo } from "react";
import { usePolling } from "./usePolling";
import { getHouseRankings, getStudentRankings, getStats } from "@/api/rankings";
import type { RankedStudent, HouseRanking, Stats } from "@/types";

const POLL_INTERVAL = 5000;

export function useRankings() {
  const prevStudentRanks = useRef<Map<number, number>>(new Map());

  const { data: houses, refresh: refreshHouses } = usePolling(getHouseRankings, POLL_INTERVAL);
  const { data: rawStudents, refresh: refreshStudents } = usePolling(getStudentRankings, POLL_INTERVAL);
  const { data: stats, refresh: refreshStats } = usePolling(getStats, POLL_INTERVAL);

  const students = useMemo<RankedStudent[]>(() => {
    if (!rawStudents) return [];

    const prev = prevStudentRanks.current;
    const withChanges = rawStudents.map((s) => {
      const prevRank = prev.get(s.id);
      const rankChange = prevRank !== undefined ? prevRank - s.rank : 0;
      return { ...s, rankChange };
    });

    // Update prev ranks for next comparison
    const newRanks = new Map<number, number>();
    rawStudents.forEach((s) => newRanks.set(s.id, s.rank));
    prevStudentRanks.current = newRanks;

    return withChanges;
  }, [rawStudents]);

  const refresh = () => {
    refreshHouses();
    refreshStudents();
    refreshStats();
  };

  return {
    houses: houses ?? ([] as HouseRanking[]),
    students,
    stats: stats ?? ({ total_laps: 0, total_participants: 0 } as Stats),
    refresh,
  };
}
