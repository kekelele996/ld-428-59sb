import { useEffect, useState } from 'react';

import { sessionOptions } from '../../constants/sessionOptions';
import { useReservationStore } from '../../stores/reservationStore';
import { ExhibitionStatus, VisitSession } from '../../types/enums';
import type { Exhibition } from '../../types/exhibition';

function todayString(): string {
  const now = new Date();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  const day = String(now.getDate()).padStart(2, '0');
  return `${now.getFullYear()}-${month}-${day}`;
}

// ReservationPanel 展览预约入口：选日期与场次、填人数和手机号后提交。
export function ReservationPanel({ exhibition }: { exhibition: Exhibition }) {
  const { availability, submitting, loadAvailability, submitReservation } = useReservationStore();
  const today = todayString();
  const minDate = exhibition.startDate > today ? exhibition.startDate : today;
  const [visitDate, setVisitDate] = useState(minDate);
  const [session, setSession] = useState<VisitSession>(VisitSession.Morning);
  const [visitorCount, setVisitorCount] = useState(1);
  const [phone, setPhone] = useState('');
  const [feedback, setFeedback] = useState<{ kind: 'success' | 'error'; text: string } | null>(null);

  // 展览未发布或已结束时停止预约。
  const reservable = exhibition.status === ExhibitionStatus.Active && exhibition.endDate >= today;

  useEffect(() => {
    if (reservable && visitDate) {
      void loadAvailability(exhibition.id, visitDate);
    }
  }, [reservable, exhibition.id, visitDate, loadAvailability]);

  if (!reservable) {
    return (
      <section className="border border-ink/15 bg-rice/70 p-6">
        <h3 className="font-display text-2xl">预约参观</h3>
        <p className="mt-3 text-sm leading-6 text-ink/60">展览未发布或已结束，暂不接受预约。</p>
      </section>
    );
  }

  const remainingOf = (value: VisitSession) => availability.find((item) => item.session === value)?.remaining;

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setFeedback(null);
    try {
      const result = await submitReservation(exhibition.id, { visitDate, session, visitorCount, phone });
      const sessionText = sessionOptions.find((option) => option.value === session)?.label ?? session;
      setFeedback({
        kind: 'success',
        text: result.duplicated
          ? '该手机号已预约过此场次，不会重复占用名额。'
          : `预约成功：${visitDate} ${sessionText}，共 ${result.reservation.visitorCount} 人，请按时到场。`,
      });
      void loadAvailability(exhibition.id, visitDate);
    } catch (error) {
      setFeedback({ kind: 'error', text: error instanceof Error ? error.message : '预约失败，请稍后再试。' });
    }
  };

  return (
    <section className="border border-ink/15 bg-rice/70 p-6">
      <h3 className="font-display text-2xl">预约参观</h3>
      <p className="mt-2 text-sm text-ink/60">选择日期与场次，填写人数和手机号即可提交，无需另行联系馆方。</p>
      <form className="mt-5 space-y-4" onSubmit={handleSubmit}>
        <label className="block text-sm">
          <span className="text-ink/65">参观日期</span>
          <input
            type="date"
            required
            min={minDate}
            max={exhibition.endDate}
            value={visitDate}
            onChange={(event) => setVisitDate(event.target.value)}
            className="mt-1 w-full border border-ink/20 bg-rice px-4 py-3 outline-none focus:border-clay"
          />
        </label>
        <div className="text-sm">
          <span className="text-ink/65">场次（每场限 {availability[0]?.capacity ?? 50} 人）</span>
          <div className="mt-2 grid gap-2">
            {sessionOptions.map((option) => {
              const remaining = remainingOf(option.value);
              const full = remaining !== undefined && remaining <= 0;
              return (
                <button
                  key={option.value}
                  type="button"
                  disabled={full}
                  onClick={() => setSession(option.value)}
                  className={`flex items-center justify-between border px-4 py-3 text-left text-sm transition ${
                    session === option.value ? 'border-ink bg-ink text-rice' : 'border-ink/20 hover:border-ink/50'
                  } ${full ? 'cursor-not-allowed opacity-40' : ''}`}
                >
                  <span>
                    {option.label} <span className="text-xs opacity-70">{option.time}</span>
                  </span>
                  <span className="text-xs">{remaining === undefined ? '余票查询中' : full ? '已满' : `剩余 ${remaining} 名`}</span>
                </button>
              );
            })}
          </div>
        </div>
        <div className="grid grid-cols-2 gap-3">
          <label className="block text-sm">
            <span className="text-ink/65">参观人数</span>
            <input
              type="number"
              required
              min={1}
              max={10}
              value={visitorCount}
              onChange={(event) => setVisitorCount(Number(event.target.value))}
              className="mt-1 w-full border border-ink/20 bg-rice px-4 py-3 outline-none focus:border-clay"
            />
          </label>
          <label className="block text-sm">
            <span className="text-ink/65">手机号</span>
            <input
              type="tel"
              required
              pattern="1[3-9][0-9]{9}"
              placeholder="用于接收通知"
              value={phone}
              onChange={(event) => setPhone(event.target.value)}
              className="mt-1 w-full border border-ink/20 bg-rice px-4 py-3 outline-none focus:border-clay"
            />
          </label>
        </div>
        <button
          type="submit"
          disabled={submitting}
          className="w-full bg-ink px-5 py-3 text-sm font-semibold text-rice hover:bg-clay disabled:opacity-50"
        >
          {submitting ? '提交中…' : '提交预约'}
        </button>
        {feedback && (
          <p className={`text-sm ${feedback.kind === 'success' ? 'text-emerald-700' : 'text-red-700'}`}>{feedback.text}</p>
        )}
      </form>
    </section>
  );
}
