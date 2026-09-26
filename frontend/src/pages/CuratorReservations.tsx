import { useEffect, useMemo, useState } from 'react';

import { EmptyState } from '../components/common/EmptyState';
import { StatusBadge } from '../components/common/StatusBadge';
import { useExhibitionStore } from '../stores/exhibitionStore';
import { useReservationStore } from '../stores/reservationStore';
import { useSessionStore } from '../stores/sessionStore';
import type { Reservation, Session } from '../types';
import { ReservationStatus } from '../types';
import { ApiError } from '../utils/request';
import { SessionForm } from '../components/reservation/SessionForm';

function formatTime(value: string): string {
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString('zh-CN', { hour12: false });
}

function ReservationRow({
  reservation,
  sessions,
  onCancel,
}: {
  reservation: Reservation;
  sessions: Session[];
  onCancel: (id: string) => Promise<void>;
}) {
  const session = sessions.find((item) => item.id === reservation.sessionId);
  const booked = reservation.status === ReservationStatus.Booked;
  return (
    <tr className="border-b border-ink/10">
      <td className="px-4 py-3 text-sm">{session ? `${session.date} ${session.startTime}-${session.endTime}` : reservation.sessionId}</td>
      <td className="px-4 py-3 text-sm">{reservation.phone}</td>
      <td className="px-4 py-3 text-sm">{reservation.partySize} 人</td>
      <td className="px-4 py-3 text-sm">{formatTime(reservation.createdAt)}</td>
      <td className="px-4 py-3">
        <StatusBadge status={booked ? 'Active' : 'Archived'} />
        <span className="ml-2 text-xs text-ink/50">{booked ? '有效' : '已取消'}</span>
      </td>
      <td className="px-4 py-3 text-right">
        {booked && (
          <button
            onClick={() => void onCancel(reservation.id)}
            className="border border-clay px-3 py-1.5 text-xs font-semibold text-clay transition hover:bg-clay hover:text-rice"
          >
            取消预约
          </button>
        )}
      </td>
    </tr>
  );
}

export function CuratorReservations() {
  const { exhibitions, loadExhibitions } = useExhibitionStore();
  const { sessions, loadSessions } = useSessionStore();
  const { reservations, loadReservations, cancel } = useReservationStore();
  const [selectedExhibition, setSelectedExhibition] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    void loadExhibitions();
  }, [loadExhibitions]);

  useEffect(() => {
    if (exhibitions.length > 0 && !selectedExhibition) setSelectedExhibition(exhibitions[0].id);
  }, [exhibitions, selectedExhibition]);

  useEffect(() => {
    if (selectedExhibition) {
      setMessage('');
      void loadSessions(selectedExhibition).catch(() => undefined);
      void loadReservations(selectedExhibition).catch((err) => {
        setMessage(err instanceof ApiError ? err.message : '加载预约记录失败。');
      });
    }
  }, [selectedExhibition, loadSessions, loadReservations]);

  const exhibition = exhibitions.find((item) => item.id === selectedExhibition);
  const bookedCount = useMemo(
    () => reservations.filter((item) => item.status === ReservationStatus.Booked).length,
    [reservations],
  );
  const seatsHeld = useMemo(
    () => reservations.filter((item) => item.status === ReservationStatus.Booked).reduce((sum, item) => sum + item.partySize, 0),
    [reservations],
  );

  const handleCancel = async (id: string) => {
    try {
      await cancel(id);
      await loadSessions(selectedExhibition);
      setMessage('已取消，名额立即释放。');
    } catch (err) {
      setMessage(err instanceof ApiError ? err.message : '取消失败，请稍后再试。');
    }
  };

  return (
    <main className="page-shell">
      <div className="mx-auto max-w-6xl px-6 py-8">
        <p className="text-sm uppercase tracking-[0.25em] text-clay">Curator desk</p>
        <h1 className="mt-2 font-display text-5xl">预约管理</h1>

        <div className="mt-6 flex flex-wrap items-center gap-3">
          <label className="text-sm text-ink/60">选择展览</label>
          <select
            value={selectedExhibition}
            onChange={(event) => setSelectedExhibition(event.target.value)}
            className="border border-ink/20 bg-white px-4 py-2.5 text-sm outline-none focus:border-clay"
          >
            {exhibitions.map((item) => (
              <option key={item.id} value={item.id}>{item.title}（{item.status}）</option>
            ))}
          </select>
          {exhibition && <StatusBadge status={exhibition.status} />}
        </div>

        {message && <p className="mt-4 text-sm text-clay">{message}</p>}

        <section className="mt-8 grid gap-4 md:grid-cols-3">
          <div className="border border-ink/15 p-5">
            <p className="text-xs uppercase tracking-wider text-ink/45">有效预约</p>
            <p className="mt-2 font-display text-4xl">{bookedCount}</p>
          </div>
          <div className="border border-ink/15 p-5">
            <p className="text-xs uppercase tracking-wider text-ink/45">占用名额</p>
            <p className="mt-2 font-display text-4xl">{seatsHeld}</p>
          </div>
          <div className="border border-ink/15 p-5">
            <p className="text-xs uppercase tracking-wider text-ink/45">场次数</p>
            <p className="mt-2 font-display text-4xl">{sessions.length}</p>
          </div>
        </section>

        {exhibition && (
          <section className="mt-8">
            <h2 className="font-display text-2xl">场次安排</h2>
            <div className="mt-4 grid gap-3 md:grid-cols-2">
              {sessions.map((session) => (
                <div key={session.id} className="flex items-center justify-between border border-ink/15 bg-white px-4 py-3 text-sm">
                  <span>
                    <span className="font-semibold">{session.date}</span>
                    <span className="ml-2 text-ink/60">{session.startTime} - {session.endTime}</span>
                  </span>
                  <span className={session.bookedCount >= session.capacity ? 'font-semibold text-clay' : 'text-ink/55'}>
                    {session.bookedCount} / {session.capacity}
                  </span>
                </div>
              ))}
            </div>
            <SessionForm
              exhibitionId={selectedExhibition}
              startDate={exhibition.startDate}
              endDate={exhibition.endDate}
              onCreated={() => void loadSessions(selectedExhibition)}
            />
          </section>
        )}

        <section className="mt-10">
          <h2 className="font-display text-2xl">预约记录</h2>
          {reservations.length === 0 ? (
            <EmptyState title="暂无预约记录" description="观众在展览详情页提交预约后会出现在这里。" />
          ) : (
            <div className="mt-4 overflow-x-auto border border-ink/15 bg-white">
              <table className="w-full">
                <thead className="bg-rice text-left text-xs uppercase tracking-wider text-ink/50">
                  <tr>
                    <th className="px-4 py-3">场次</th>
                    <th className="px-4 py-3">手机号</th>
                    <th className="px-4 py-3">人数</th>
                    <th className="px-4 py-3">提交时间</th>
                    <th className="px-4 py-3">状态</th>
                    <th className="px-4 py-3" />
                  </tr>
                </thead>
                <tbody>
                  {reservations.map((reservation) => (
                    <ReservationRow key={reservation.id} reservation={reservation} sessions={sessions} onCancel={handleCancel} />
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
