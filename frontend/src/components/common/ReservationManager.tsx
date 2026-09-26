import { useEffect, useState } from 'react';

import { sessionLabel } from '../../constants/sessionOptions';
import { useExhibitionStore } from '../../stores/exhibitionStore';
import { useReservationStore } from '../../stores/reservationStore';
import { ReservationStatus } from '../../types/enums';
import { EmptyState } from './EmptyState';

// ReservationManager 策展人在工作台查看并取消观众预约；取消后名额立即释放。
export function ReservationManager() {
  const { exhibitions, loadExhibitions } = useExhibitionStore();
  const { reservations, loadReservations, cancel } = useReservationStore();
  const [exhibitionId, setExhibitionId] = useState('');
  const [cancellingId, setCancellingId] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    void loadExhibitions();
  }, [loadExhibitions]);

  useEffect(() => {
    if (!exhibitionId && exhibitions.length > 0) {
      setExhibitionId(exhibitions[0].id);
    }
  }, [exhibitions, exhibitionId]);

  useEffect(() => {
    if (exhibitionId) {
      void loadReservations(exhibitionId);
    }
  }, [exhibitionId, loadReservations]);

  const handleCancel = async (id: string) => {
    setError('');
    setCancellingId(id);
    try {
      await cancel(id);
    } catch (err) {
      setError(err instanceof Error ? err.message : '取消失败，请稍后再试。');
    } finally {
      setCancellingId('');
    }
  };

  return (
    <section className="mt-12">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm uppercase tracking-[0.25em] text-clay">Reservations</p>
          <h2 className="mt-2 font-display text-4xl">预约管理</h2>
        </div>
        <select
          value={exhibitionId}
          onChange={(event) => setExhibitionId(event.target.value)}
          className="border border-ink/20 bg-rice px-4 py-3 text-sm outline-none focus:border-clay"
        >
          {exhibitions.map((exhibition) => (
            <option key={exhibition.id} value={exhibition.id}>{exhibition.title}</option>
          ))}
        </select>
      </div>
      {error && <p className="mt-4 text-sm text-red-700">{error}</p>}
      <div className="mt-6">
        {reservations.length === 0 ? (
          <EmptyState title="暂无预约记录" description="观众在展览详情页提交预约后会出现在这里。" />
        ) : (
          <div className="overflow-x-auto border border-ink/15">
            <table className="w-full min-w-[720px] text-left text-sm">
              <thead className="bg-ink text-rice">
                <tr>
                  <th className="px-4 py-3 font-medium">参观日期</th>
                  <th className="px-4 py-3 font-medium">场次</th>
                  <th className="px-4 py-3 font-medium">人数</th>
                  <th className="px-4 py-3 font-medium">手机号</th>
                  <th className="px-4 py-3 font-medium">状态</th>
                  <th className="px-4 py-3 font-medium">提交时间</th>
                  <th className="px-4 py-3 font-medium">操作</th>
                </tr>
              </thead>
              <tbody>
                {reservations.map((reservation) => (
                  <tr key={reservation.id} className="border-t border-ink/10">
                    <td className="px-4 py-3">{reservation.visitDate}</td>
                    <td className="px-4 py-3">{sessionLabel(reservation.session)}</td>
                    <td className="px-4 py-3">{reservation.visitorCount}</td>
                    <td className="px-4 py-3">{reservation.phone}</td>
                    <td className="px-4 py-3">
                      <span className={reservation.status === ReservationStatus.Confirmed ? 'text-emerald-700' : 'text-ink/45'}>
                        {reservation.status === ReservationStatus.Confirmed ? '已确认' : '已取消'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-ink/60">{new Date(reservation.createdAt).toLocaleString('zh-CN')}</td>
                    <td className="px-4 py-3">
                      {reservation.status === ReservationStatus.Confirmed && (
                        <button
                          onClick={() => void handleCancel(reservation.id)}
                          disabled={cancellingId === reservation.id}
                          className="border border-ink/25 px-3 py-1 text-xs hover:border-clay hover:text-clay disabled:opacity-50"
                        >
                          {cancellingId === reservation.id ? '取消中…' : '取消预约'}
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </section>
  );
}
