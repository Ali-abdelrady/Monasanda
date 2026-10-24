import { ChevronDown, Filter, LayoutGrid, List, Search } from "lucide-react";
import type { StatusFilter, ViewMode } from "../types";

interface Props {
  search: string;
  onSearch: (v: string) => void;
  status: StatusFilter;
  onStatus: (v: StatusFilter) => void;
  view: ViewMode;
  onView: (v: ViewMode) => void;
}

export function AttendanceToolbar({ search, onSearch, status, onStatus, view, onView }: Props) {
  return (
    <div className="toolbar">
      <div className="search-box">
        <Search size={18} />
        <input
          value={search}
          onChange={(e) => onSearch(e.target.value)}
          placeholder="ابحث بالاسم أو الرقم العسكري أو التليفون..."
        />
        <kbd>⌘ K</kbd>
      </div>
      <div className="toolbar-controls">
        <div className="filter-select">
          <Filter size={15} />
          <select value={status} onChange={(e) => onStatus(e.target.value as StatusFilter)}>
            <option value="all">كل الحالات</option>
            <option value="present">حضروا</option>
            <option value="absent">لم يحضروا</option>
          </select>
          <ChevronDown size={14} />
        </div>
        <div className="view-toggle">
          <button
            className={view === "table" ? "selected" : ""}
            onClick={() => onView("table")}
            title="عرض الجدول"
          >
            <List size={18} />
          </button>
          <button
            className={view === "cards" ? "selected" : ""}
            onClick={() => onView("cards")}
            title="عرض الكروت"
          >
            <LayoutGrid size={17} />
          </button>
        </div>
      </div>
    </div>
  );
}
