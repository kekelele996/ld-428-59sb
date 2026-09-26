import { VisitSession } from '../types/enums';

export const sessionOptions = [
  { value: VisitSession.Morning, label: '上午场', time: '09:00 - 12:00' },
  { value: VisitSession.Afternoon, label: '下午场', time: '13:00 - 17:00' },
  { value: VisitSession.Evening, label: '晚间场', time: '18:00 - 21:00' },
];

export function sessionLabel(session: VisitSession): string {
  return sessionOptions.find((option) => option.value === session)?.label ?? session;
}
