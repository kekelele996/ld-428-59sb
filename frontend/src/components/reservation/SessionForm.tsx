import { useState } from 'react';

import { useSessionStore } from '../../stores/sessionStore';
import { ApiError } from '../../utils/request';

export function SessionForm({
  exhibitionId,
  startDate,
  endDate,
  onCreated,
}: {
  exhibitionId: string;
  startDate: string;
  endDate: string;
  onCreated: () => void;
}) {
  const { addSession } = useSessionStore();
  const [open, setOpen] = useState(false);
  const [date, setDate] = useState(startDate ?? '');
  const [startTime, setStartTime] = useState('10:00');
  const [endTime, setEndTime] = useState('12:00');
  const [capacity, setCapacity] = useState(20);
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError('');
    setSaving(true);
    try {
      await addSession(exhibitionId, { date, startTime, endTime, capacity });
      setOpen(false);
      onCreated();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : '创建场次失败。');
    } finally {
      setSaving(false);
    }
  };

  if (!open) {
    return (
      <button
        onClick={() => setOpen(true)}
        className="mt-4 border border-ink px-5 py-2.5 text-sm font-semibold hover:bg-ink hover:text-rice"
      >
        + 新增场次
      </button>
    );
  }

  return (
    <form className="mt-4 grid gap-3 border border-ink/15 bg-rice/60 p-5 md:grid-cols-5" onSubmit={submit}>
      <label className="text-xs">
        <span className="text-ink/55">日期（展期 {startDate} ~ {endDate}）</span>
        <input
          type="date"
          required
          min={startDate}
          max={endDate}
          value={date}
          onChange={(event) => setDate(event.target.value)}
          className="mt-1 w-full border border-ink/20 bg-white px-3 py-2 text-sm outline-none focus:border-clay"
        />
      </label>
      <label className="text-xs">
        <span className="text-ink/55">开始时间</span>
        <input
          type="time"
          required
          value={startTime}
          onChange={(event) => setStartTime(event.target.value)}
          className="mt-1 w-full border border-ink/20 bg-white px-3 py-2 text-sm outline-none focus:border-clay"
        />
      </label>
      <label className="text-xs">
        <span className="text-ink/55">结束时间</span>
        <input
          type="time"
          required
          value={endTime}
          onChange={(event) => setEndTime(event.target.value)}
          className="mt-1 w-full border border-ink/20 bg-white px-3 py-2 text-sm outline-none focus:border-clay"
        />
      </label>
      <label className="text-xs">
        <span className="text-ink/55">人数上限</span>
        <input
          type="number"
          required
          min={1}
          max={10000}
          value={capacity}
          onChange={(event) => setCapacity(Number(event.target.value))}
          className="mt-1 w-full border border-ink/20 bg-white px-3 py-2 text-sm outline-none focus:border-clay"
        />
      </label>
      <div className="flex items-end gap-2">
        <button type="submit" disabled={saving} className="bg-ink px-4 py-2 text-sm font-semibold text-rice hover:bg-clay disabled:opacity-40">
          {saving ? '保存中…' : '保存'}
        </button>
        <button type="button" onClick={() => setOpen(false)} className="border border-ink/20 px-4 py-2 text-sm hover:border-ink">
          取消
        </button>
      </div>
      {error && <p className="text-sm text-clay md:col-span-5">{error}</p>}
    </form>
  );
}
