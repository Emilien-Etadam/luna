type EventModel = {
  id: string;
  calendar: string;
  name: string;
  desc: string;
  color: string;
  date: {
    start: Date;
    end: Date;
    allDay: boolean;
    recurrence: any;
  };
  overridden: boolean;
  can_edit: boolean;
  can_delete: boolean;
  participant_colors?: string[];
  /** Après fusion multi-calendriers (vue publique) : noms des propriétaires concernés. */
  calendar_owner_names?: string[];
};

type EventModelChanges = {
  name: boolean;
  desc: boolean;
  color: boolean;
  date: boolean;
}