import { useEffect, useMemo, useState } from 'react';

import { useExhibitionStore } from '../../stores/exhibitionStore';
import { useSessionStore } from '../../stores/sessionStore';
import type { Exhibition, Session } from '../../types';
import { ExhibitionStatus } from '../../types/enums';
import { createReservation } from '../../api/reservation';
import { ApiError } from '../../utils/request';

const MAX_PARTY = 10;

// isReservationOpen 展览仅在“已发布且在展”时接受预约。
function isReservationOpen(exhibition: Exhibition): boolean {
  return exhibition.reviewStatus === 'Approved' && exhibition.status === ExhibitionStatus.Active;
}

function closedReason(exhibition: Exhibition): string {
  if (exhibition.reviewStatus !== 'Approved' || exhibition.status === ExhibitionStatus.Planning) {
    return '展览尚未发布，预约通道暂未开放。';
  }
  if (exhibition.status === ExhibitionStatus.Ended || exhibition.status === ExhibitionStatus.Archived) {
    return '展览已结束，不再接受预约。';
  }
  return '当前无法预约。';
}

function remaining(session: Session): number {
  return Math.max(0, session.capacity - session.bookedCount);
}

export function ReservationPanel({ exhibitionId }: { exhibitionId: string }) {
  const { exhibitions } = useExhibitionStore();
  const { sessions, loadSessions } = useSessionStore();
  const exhibition = exhibitions.find((item) => item.id === exhibitionId);

  const [sessionId, setSessionId] = useState('');
  const [phone, setPhone] = useState('');
  const [partySize, setPartySize] = useState(1);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    void loadSessions(exhibitionId);
  }, [exhibitionId, loadSessions]);

  const selectableSessions = useMemo(() => sessions.filter((session) => remaining(session) > 0), [sessions]);

  useEffect(() => {
    if (!sessionId && selectableSessions.length > 0) setSessionId(selectableSessions[0].id);
  }, [selectableSessions, sessionId]);

  if (!exhibition) return null;

  if (!isReservationOpen(exhibition)) {
    return (
      <section className="border border-ink/15 bg-rice/60 p-6">
        <h3 className="font-display text-2xl">预约参观</h3>
        <p className="mt-3 text-sm leading-7 text-ink/60">{closedReason(exhibition)}</p>
      </section>
    );
  }

  const selected = sessions.find((session) => session.id === sessionId);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError('');
    setSuccess(null);
    if (!sessionId) {
      setError('请选择场次。');
      return;
    }
    if (!/^1[3-9]\d{9}$/.test(phone)) {
      setError('请输入有效的 11 位手机号。');
      return;
    }
    if (selected && partySize > remaining(selected)) {
      setError(`该场次仅剩 ${remaining(selected)} 个名额。`);
      return;
    }
    setSubmitting(true);
    try {
      const reservation = await createReservation({ sessionId, phone, partySize });
      setSuccess(`预约成功！预约号 ${reservation.id}，请留意手机通知。`);
      setPhone('');
      setPartySize(1);
      await loadSessions(exhibitionId);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : '提交失败，请稍后再试。');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="border border-ink/15 bg-rice/60 p-6">
      <div className="flex items-baseline justify-between">
        <h3 className="font-display text-2xl">预约参观</h3>
        <span className="text-xs uppercase tracking-[0.2em] text-moss">Booking open</span>
      </div>
      {sessions.length === 0 ? (
        <p className="mt-3 text-sm leading-7 text-ink/60">馆方暂未开放可预约场次，请稍后再来。</p>
      ) : (
        <form className="mt-5 space-y-4" onSubmit={handleSubmit}>
          <div>
            <label className="text-xs font-semibold uppercase tracking-wider text-ink/55">选择场次</label>
            <div className="mt-2 space-y-2">
              {sessions.map((session) => {
                const left = remaining(session);
                const full = left === 0;
                return (
                  <label
                    key={session.id}
                    className={`flex cursor-pointer items-center justify-between border px-4 py-3 text-sm transition ${
                      sessionId === session.id ? 'border-clay bg-white' : 'border-ink/15 bg-white/50 hover:border-ink/40'
                    } ${full ? 'cursor-not-allowed opacity-50' : ''}`}
                  >
                    <span className="flex items-center gap-3">
                      <input
                        type="radio"
                        name="session"
                        value={session.id}
                        checked={sessionId === session.id}
                        disabled={full}
                        onChange={() => setSessionId(session.id)}
                      />
                      <span>
                        <span className="font-semibold">{session.date}</span>
                        <span className="ml-2 text-ink/60">{session.startTime} - {session.endTime}</span>
                      </span>
                    </span>
                    <span className={full ? 'text-clay' : 'text-ink/55'}>{full ? '已满' : `剩 ${left} 席`}</span>
                  </label>
                );
              })}
            </div>
          </div>
          <div className="grid grid-cols-[1fr_120px] gap-3">
            <div>
              <label className="text-xs font-semibold uppercase tracking-wider text-ink/55">手机号</label>
              <input
                value={phone}
                onChange={(event) => setPhone(event.target.value.replace(/\D/g, '').slice(0, 11))}
                inputMode="numeric"
                placeholder="用于接收参观通知"
                className="mt-2 w-full border border-ink/20 bg-white px-4 py-3 text-sm outline-none focus:border-clay"
              />
            </div>
            <div>
              <label className="text-xs font-semibold uppercase tracking-wider text-ink/55">人数</label>
              <select
                value={partySize}
                onChange={(event) => setPartySize(Number(event.target.value))}
                className="mt-2 w-full border border-ink/20 bg-white px-3 py-3 text-sm outline-none focus:border-clay"
              >
                {Array.from({ length: MAX_PARTY }, (_, index) => index + 1).map((n) => (
                  <option key={n} value={n}>{n} 人</option>
                ))}
              </select>
            </div>
          </div>
          {error && <p className="text-sm text-clay">{error}</p>}
          {success && <p className="text-sm text-moss">{success}</p>}
          <button
            type="submit"
            disabled={submitting || selectableSessions.length === 0}
            className="w-full bg-ink px-5 py-3 text-sm font-semibold text-rice transition hover:bg-clay disabled:cursor-not-allowed disabled:opacity-40"
          >
            {submitting ? '提交中…' : '提交预约'}
          </button>
          <p className="text-xs leading-5 text-ink/45">每个场次有人数上限，约满即止；同一手机号在同一场次仅可预约一次。</p>
        </form>
      )}
    </section>
  );
}
