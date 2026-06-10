<script lang="ts">
  import { Github } from "svelte-simples"
  import { CalendarDays, CalendarRange, Calendar as CalendarSingle, Copyleft, List, PlusIcon, RefreshCw, Settings, WifiOff } from "lucide-svelte";
  import { onMount, setContext, untrack } from "svelte";

  import Calendar from "../../components/calendar/Calendar.svelte";
  import CalendarEntry from "../../components/calendar/CalendarEntry.svelte";
  import CalendarModal from "../../components/modals/CalendarModal.svelte";
  import EventModal from "../../components/modals/EventModal.svelte";
  import Horizontal from "../../components/layout/Horizontal.svelte";
  import IconButton from "../../components/interactive/IconButton.svelte";
  import MonthSelection from "../../components/interactive/MonthSelection.svelte";
  import SourceEntry from "../../components/calendar/SourceEntry.svelte";
  import SourceModal from "../../components/modals/SourceModal.svelte";

  import { afterNavigate, beforeNavigate } from "$app/navigation";
  import { browser } from "$app/environment";

  import SmallCalendar from "../../components/interactive/SmallCalendar.svelte";
  import { NoOp } from "$lib/client/placeholders";
  import { getMetadata } from "$lib/client/data/metadata.svelte";
  import { getRepository } from "$lib/client/data/repository.svelte";
  import { queueNotification } from "$lib/client/notifications";
  import { getConnectivity, Reachability } from "$lib/client/data/connectivity.svelte";
  import Button from "../../components/interactive/Button.svelte";
  import DayViewModal from "../../components/modals/DayViewModal.svelte";
  import { getDayIndex, getLocalDayKey, getMsUntilNextMidnight, isInRange, isSameDay } from "../../lib/common/date";
  import { compareEventsByStartDate } from "../../lib/common/comparators";
  import { GetEventColor } from "$lib/common/colors";
  import SourceWizardModal from "../../components/modals/SourceWizardModal.svelte";
  import SettingsModal from "../../components/modals/SettingsModal.svelte";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import ThemeToggle from "../../components/interactive/ThemeToggle.svelte";
  import { ColorKeys } from "../../types/colors";
  import { page } from "$app/state";
  import CreditsModal from "../../components/modals/CreditsModal.svelte";

  /* Singletons */
  const settings = getSettings();
  const metadata = getMetadata();
  const repository = getRepository();
  const connectivity = getConnectivity();
  const publicReadonly = $derived(!settings.userData.id);

  /* Constants */
  let autoRefreshInterval = 1000 * 60; // 1 minute

  /* View logic */
  let view: "month" | "week" | "day" | "agenda" = $derived.by(() => {
    const stored = page.url.searchParams.get("view");
    if (!stored || !["month", "week", "day", "agenda"].includes(stored)) return "month"
    return stored as "month" | "week" | "day" | "agenda";
  });

  let today = $state(new Date());
  let date = $derived.by(() => {
    const stored = page.url.searchParams.get("date");
    if (!stored) return new Date(today);
    const parsed = new Date(stored);
    return parsed;
  });

  $effect(() => {
    const url = new URL(window.location.toString());

    url.searchParams.set("view", view);
    url.searchParams.set("date", date.toISOString().split("T")[0]);

    history.replaceState(history.state, '', url);
  })

  function getVisibleRange(date: Date, view: "month" | "week" | "day" | "agenda"): { start: Date, end: Date } {
    const rangeStart = new Date(date);
    const rangeEnd = new Date(date);
    rangeStart.setHours(0, 0, 0, 0);
    rangeEnd.setHours(23, 59, 59, 999);
    switch (view) {
      case "month":
        rangeStart.setDate(1);
        rangeStart.setDate(rangeStart.getDate() - getDayIndex(rangeStart));
        rangeEnd.setMonth(rangeEnd.getMonth() + 1);
        rangeEnd.setDate(0);
        rangeEnd.setDate(rangeEnd.getDate() + 6 - getDayIndex(rangeEnd));
        break;
      case "week":
        rangeStart.setDate(date.getDate() - getDayIndex(rangeStart));
        rangeEnd.setDate(rangeStart.getDate() + 7);
        break;
      case "day":
      case "agenda":
        rangeStart.setDate(1);
        rangeEnd.setMonth(rangeEnd.getMonth() + 1);
        rangeEnd.setDate(0);
        break;
      default:
    }
    return { start: rangeStart, end: rangeEnd };
  }
  let todayInRange = $derived.by(() => {
    const range = getVisibleRange(date, view);
    return isInRange(today, range.start, range.end)
  });

  function seeToday() {
    today = new Date();
    date = new Date(today);
  }

  function isAnchoredOnToday(anchorDate: Date, currentView: typeof view, referenceToday: Date): boolean {
    if (currentView === "week") {
      const range = getVisibleRange(anchorDate, "week");
      return isInRange(referenceToday, range.start, range.end);
    }
    if (currentView === "day" || currentView === "agenda") {
      return isSameDay(anchorDate, referenceToday);
    }
    return isSameDay(anchorDate, referenceToday);
  }

  function syncTodayIfNeeded() {
    const now = new Date();
    if (isSameDay(today, now)) return;

    const shouldAdvanceDate = isAnchoredOnToday(date, view, today);
    today = now;
    if (shouldAdvanceDate) {
      date = new Date(today);
    }
  }

  function goToView(target: typeof view) {
    syncTodayIfNeeded();
    if (target === "day" || target === "week" || target === "agenda") {
      date = new Date();
    }
    view = target;
  }

  let dayChangeTimeout: ReturnType<typeof setTimeout> | undefined;

  function scheduleDayChangeCheck() {
    if (!browser) return;
    clearTimeout(dayChangeTimeout);
    dayChangeTimeout = setTimeout(() => {
      syncTodayIfNeeded();
      scheduleDayChangeCheck();
    }, getMsUntilNextMidnight());
  }

  function handleVisibilityChange() {
    if (document.visibilityState === "visible") {
      syncTodayIfNeeded();
    }
  }

  function smallCalendarClick(clickedDate: Date) {
    const range = getVisibleRange(date, view);
    if (isInRange(clickedDate, range.start, range.end) && clickedDate.getMonth() === date.getMonth()) {
      showDateModal(clickedDate, repository.events
        .filter((event) => event.date.start.getTime() <= clickedDate.getTime() + 24 * 60 * 60 * 1000 && event.date.end.getTime() >= clickedDate.getTime())
        .sort(compareEventsByStartDate)
      );
    } else {
      date = clickedDate;
    }
  }

  /* Fetching logic */
  let isLoading: boolean = $derived(getMetadata().loadingData);
  let loaderAnimation = $state(false);
  $effect(() => {
    if (isLoading) loaderAnimation = true;
  })

  afterNavigate(() => {
    void refresh();
  });

  beforeNavigate((args) => {
    if (args.to === null) return;
    clearTimeout(spooledRefresh);
    spooledRefresh = undefined; 
  });

  let spooledRefresh: (ReturnType<typeof setTimeout> | undefined) = $state(undefined);
  let refreshRequestId = 0;
  let manualRefreshInFlight = false;

  function formatFetchError(err: unknown): string {
    const message = err instanceof Error ? err.message : String(err);
    if (message.startsWith("Failed to fetch events")) return message;
    return `Failed to fetch events: ${message}`;
  }

  function refresh(force = false): Promise<void> {
    if (browser) syncTodayIfNeeded();
    const range = getVisibleRange(date, view);
    const requestId = ++refreshRequestId;

    const promise = getRepository().getAllEvents(range.start, range.end, force).catch((err) => {
      if (requestId !== refreshRequestId) return;
      queueNotification(ColorKeys.Danger, formatFetchError(err));
    }).then(() => {});

    connectivity.check();

    if (!publicReadonly) {
      clearTimeout(spooledRefresh);
      spooledRefresh = setTimeout(() => {
        void refresh();
      }, autoRefreshInterval);
    }

    return promise;
  }

  function forceRefresh() {
    if (manualRefreshInFlight || isLoading) return;
    manualRefreshInFlight = true;
    getRepository().invalidateCache();
    void refresh(true).finally(() => {
      manualRefreshInFlight = false;
    });
  }

  $effect(() => {
    ((date: Date, view: "month" | "week" | "day" | "agenda") => {
      untrack(() => {
        if (!browser) return;
        void refresh();
      });
    })(date, view);
  });

  onMount(() => {
    if (!browser) return;
    getRepository().getSources().catch(NoOp);
    scheduleDayChangeCheck();
    document.addEventListener("visibilitychange", handleVisibilityChange);
    return () => {
      clearTimeout(dayChangeTimeout);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    };
  });

  /* Single instance modal logic */
  let showSourceWizardModal: () => any = $state(NoOp);

  let showNewSourceModal: () => any = $state(NoOp);
  const showNewSourceModalInternal = () => { return showNewSourceModal(); };
  setContext("showNewSourceModal", showNewSourceModalInternal);

  let showSourceModal: (source: SourceModel) => any = $state(NoOp);
  const showSourceModalInternal = (source: SourceModel) => { return showSourceModal(source); };
  setContext("showSourceModal", showSourceModalInternal);

  let showNewCalendarModal: () => any = $state(NoOp);

  let showCalendarModal: (calendar: CalendarModel) => any = $state(NoOp);
  const showCalendarModalInternal = (calendar: CalendarModel) => { return showCalendarModal(calendar); };
  setContext("showCalendarModal", showCalendarModalInternal);

  let showNewEventModal: (date: Date) => any = $state(NoOp);
  const showNewEventModalInternal = (date: Date) => { return showNewEventModal(date); };
  setContext("showNewEventModal", showNewEventModalInternal);

  let showEventModal: (event: EventModel) => any = $state(NoOp);
  const showEventModalInternal = (event: EventModel) => { return showEventModal(event); };
  setContext("showEventModal", showEventModalInternal);

  let showDateModal: (date: Date, events: (EventModel | null)[]) => any = $state(NoOp);
  const showDateModalInternal = (date: Date, events: (EventModel | null)[]) => { return showDateModal(date, events); };
  setContext("showDateModal", showDateModalInternal);

  let showSettingsModal: () => any = $state(NoOp);

  let showCreditsModal: () => any = $state(NoOp);

  const timezone = $derived(Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC");
  const connectivityLabel = $derived.by(() => {
    switch (connectivity.reachable) {
      case Reachability.Database:
        return "Sync OK";
      case Reachability.Backend:
        return "DB indisponible";
      case Reachability.Frontend:
        return "Backend indisponible";
      case Reachability.None:
        return "Frontend indisponible";
      case Reachability.Incompatible:
        return "Versions incompatibles";
      default:
        return "Etat inconnu";
    }
  });
  const nextEventLabel = $derived.by(() => {
    const now = Date.now();
    const next = repository.events
      .filter((event) => event.date.end.getTime() > now)
      .sort(compareEventsByStartDate)[0];

    if (!next) return "Prochain: aucun";
    const time = next.date.start.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    return `Prochain: ${time} ${next.name}`;
  });

  const agendaDays = $derived.by(() => {
    const range = getVisibleRange(date, view);
    const start = new Date(range.start);
    start.setHours(0, 0, 0, 0);
    const end = new Date(range.end);
    end.setHours(23, 59, 59, 999);

    const todayMidnight = new Date(today);
    todayMidnight.setHours(0, 0, 0, 0);

    const isCurrentMonthShown =
      view === "agenda" &&
      date.getFullYear() === todayMidnight.getFullYear() &&
      date.getMonth() === todayMidnight.getMonth();

    const cutoff = isCurrentMonthShown ? todayMidnight : start;

    const inRangeEvents = repository.events
      .filter((event) => event.date.end.getTime() > cutoff.getTime() && event.date.start.getTime() <= end.getTime())
      .sort(compareEventsByStartDate);

    const byDay = new Map<string, { date: Date, events: EventModel[] }>();

    if (isCurrentMonthShown) {
      const key = getLocalDayKey(todayMidnight);
      byDay.set(key, { date: new Date(todayMidnight), events: [] });
    }

    for (const event of inRangeEvents) {
      const day = new Date(event.date.start);
      day.setHours(0, 0, 0, 0);
      if (day.getTime() < cutoff.getTime()) day.setTime(cutoff.getTime());
      const key = getLocalDayKey(day);
      if (!byDay.has(key)) byDay.set(key, { date: day, events: [] });
      byDay.get(key)?.events.push(event);
    }

    return Array.from(byDay.values()).sort((a, b) => a.date.getTime() - b.date.getTime());
  });

  const calendarGranularity = $derived(view === "agenda" ? "month" : view);
  const todayKey = $derived(getLocalDayKey(today));
  let agendaContainer: HTMLElement | null = $state(null);

  $effect(() => {
    if (view !== "agenda" || agendaDays.length === 0) return;
    const container = agendaContainer;
    if (!container) return;
    requestAnimationFrame(() => {
      const todayRow = container.querySelector<HTMLElement>("[data-agenda-today='true']");
      if (todayRow) {
        container.scrollTop = todayRow.offsetTop - container.offsetTop;
      } else {
        container.scrollTop = 0;
      }
    });
  });

  function openAgendaEvent(event: EventModel) {
    showEventModal(event).then((updatedEvent: EventModel) => event = updatedEvent).catch(NoOp);
  }

  function agendaItemKeydown(e: KeyboardEvent, event: EventModel) {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      openAgendaEvent(event);
    }
  }

  function agendaParticipantColors(event: EventModel): string[] {
    const colors = new Set(event.participant_colors || []);
    const own = GetEventColor(event);
    if (own) colors.add(own);
    return Array.from(colors);
  }

  function agendaIsMergedEvent(event: EventModel): boolean {
    const owners = event.calendar_owner_names?.filter(Boolean) ?? [];
    const colors = agendaParticipantColors(event);
    return owners.length > 1 || colors.length > 1;
  }

  function agendaMergedOwnersLabel(event: EventModel): string {
    if (!agendaIsMergedEvent(event)) return "";
    return (event.calendar_owner_names?.filter(Boolean) ?? []).join(", ");
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  :global(body) {
    display: flex;
    flex-direction: row;
  }

  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 0;
    grid-area: main;
    background-color: var(--bg-editor);
    overflow: hidden;
    min-width: 0;
  }
  
  div.leftPane {
    display: flex;
    flex-direction: row;
    gap: 0;
  }

  nav.activityBar {
    width: 48px;
    min-width: 48px;
    max-width: 48px;
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: 4px;
    background-color: var(--bg-activity-bar);
    border: 0;
    border-right: 1px solid var(--border-subtle);
    padding: 8px 6px;
  }

  button.activityIcon {
    all: unset;
    width: 36px;
    height: 36px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: var(--fg-muted);
    border-radius: var(--radius-2);
    box-sizing: border-box;
    position: relative;
    transition: background-color var(--transition-fast), color var(--transition-fast);
  }

  button.activityIcon:hover {
    background-color: var(--bg-hover);
    color: var(--fg-primary);
  }

  button.activityIcon:focus-visible {
    box-shadow: var(--focus-ring);
    color: var(--fg-primary);
  }

  button.activityIcon.active {
    color: var(--fg-strong);
    background-color: var(--bg-active);
  }

  /* Marqueur d'actif : barre verticale subtile à gauche */
  button.activityIcon.active::before {
    content: "";
    position: absolute;
    left: -6px;
    top: 6px;
    bottom: 6px;
    width: 2px;
    border-radius: var(--radius-pill);
    background-color: var(--accent-blue);
  }

  aside.sidebar {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
    min-width: 15em;
    width: 15em;
    max-width: 15em;
    grid-area: aside;
    background-color: var(--bg-side-bar);
    border: 0;
    border-right: 1px solid var(--border-subtle);
    padding: dimensions.$gapSmall;
  }

  div.sources {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
    overflow: auto;
    margin: 0;
    padding: dimensions.$gapTiny 0;
  }

  div.toprow {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
    justify-content: space-between;
    margin: 0;
    align-items: center;
    height: 44px;
    min-height: 44px;
    padding: 0 dimensions.$gapSmall;
    border-bottom: 1px solid var(--border-subtle);
    background-color: var(--bg-editor);
  }

  div.toprow :global(button),
  div.toprow :global(a) {
    color: var(--fg-muted);
  }

  div.toprow :global(button:hover),
  div.toprow :global(button:focus-visible),
  div.toprow :global(a:hover),
  div.toprow :global(a:focus-visible) {
    color: var(--fg-primary);
  }

  span.reachability {
    color: var(--accent-red);
    align-items: center;
    display: flex;
    flex-direction: row;
    justify-content: center;
    gap: dimensions.$gapSmaller;
    font-size: var(--font-size-sm);
  }

  span.refreshButtonWrapper {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  span.spin {
    animation: spin animations.$animationSpeedSlow animations.$cubic infinite forwards;
  }

  span.copyright {
    color: var(--fg-disabled);
    font-size: var(--font-size-xs);
    text-align: center;
    margin-top: 0;
    line-height: 1.4;
  }

  div.statusbar {
    min-height: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: dimensions.$gapSmall;
    padding: 0 dimensions.$gapSmall;
    background-color: var(--bg-status-bar);
    color: var(--fg-muted);
    font-size: var(--font-size-status-bar);
    font-weight: 400;
    border-top: 1px solid var(--border-subtle);
    white-space: nowrap;
    overflow: hidden;
  }

  div.statusbar span {
    overflow: hidden;
    text-overflow: ellipsis;
    padding: 0 dimensions.$gapSmaller;
  }

  /* ============= Agenda ============= */
  section.agenda {
    position: relative;
    height: 100%;
    overflow-y: auto;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0;
    background-color: var(--bg-editor);
  }

  div.agendaDay {
    background-color: transparent;
  }

  h3.agendaDate {
    margin: 0;
    min-height: 28px;
    box-sizing: border-box;
    padding: 4px 16px;
    font-size: var(--font-size-xs);
    font-weight: var(--font-weight-section-header);
    letter-spacing: var(--letter-spacing-section-header);
    text-transform: uppercase;
    color: var(--fg-muted);
    background-color: var(--bg-section-header);
    border-bottom: 1px solid var(--border-subtle);
    position: sticky;
    top: 0;
    z-index: 1;
  }

  div.agendaDay[data-agenda-today="true"] h3.agendaDate {
    color: var(--accent-blue);
  }

  div.agendaItem {
    display: grid;
    grid-template-columns: 92px 1fr;
    column-gap: 16px;
    align-items: center;
    box-sizing: border-box;
    min-height: 32px;
    padding: 6px 16px;
    background-color: transparent;
    border-bottom: 1px solid var(--border-subtle);
    border-left: 0;
    cursor: pointer;
    transition: background-color var(--transition-fast);
  }

  div.agendaItem:hover {
    background-color: var(--bg-hover);
  }

  div.agendaItem:focus-visible {
    outline: 0;
    box-shadow: inset 0 0 0 2px var(--accent-blue);
  }

  span.agendaTime {
    color: var(--fg-muted);
    font-family: text.$fontFamilyTime;
    font-size: var(--font-size-event-time);
    font-weight: var(--font-weight-medium);
    font-variant-numeric: tabular-nums;
    width: 92px;
    flex-shrink: 0;
  }

  span.agendaName {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
    flex: 1;
    font-size: var(--font-size-ui);
    font-weight: var(--font-weight-ui);
    color: var(--fg-primary);
    overflow: hidden;
  }

  span.agendaMainLine {
    min-width: 0;
    flex: 1;
    display: flex;
    align-items: baseline;
    gap: 6px;
    overflow: hidden;
  }

  span.agendaTitle {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  span.agendaOwners {
    flex-shrink: 0;
    max-width: min(50%, 14rem);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: var(--font-size-xs);
    color: var(--fg-muted);
    font-weight: var(--font-weight-ui);
  }

  span.eventDot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
    box-shadow: 0 0 0 1px color-mix(in srgb, currentColor 12%, transparent);
  }

  div.agendaEmpty {
    flex-grow: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--fg-muted);
    font-size: var(--font-size-sm);
    font-style: italic;
    padding: dimensions.$gapLarger;
    text-align: center;
  }

  /* ============= Responsive ============= */
  @media (max-width: 720px) {
    aside.sidebar {
      display: none;
    }
    nav.activityBar {
      width: 44px;
      min-width: 44px;
    }
    div.toprow {
      padding: 0 dimensions.$gapSmaller;
    }
    span.reachability {
      display: none;
    }
  }
</style>

{#if !publicReadonly}
  <SourceWizardModal bind:showModal={showSourceWizardModal}/>
  <SourceModal bind:showCreateModal={showNewSourceModal} bind:showModal={showSourceModal}/>
  <CalendarModal bind:showCreateModal={showNewCalendarModal} bind:showModal={showCalendarModal}/>
  <DayViewModal bind:showModal={showDateModal}/>
{/if}
<EventModal bind:showCreateModal={showNewEventModal} bind:showModal={showEventModal}/>
<SettingsModal bind:showModal={showSettingsModal} appearanceOnly={publicReadonly}/>
<CreditsModal bind:showModal={showCreditsModal}/>

<div class="leftPane">
  <nav class="activityBar" aria-label="Vues calendrier">
    <button class="activityIcon" class:active={view === "month"} onclick={() => view = "month"} title="Mois">
      <CalendarDays size={16}/>
    </button>
    <button class="activityIcon" class:active={view === "week"} onclick={() => goToView("week")} title="Semaine">
      <CalendarRange size={16}/>
    </button>
    <button class="activityIcon" class:active={view === "day"} onclick={() => goToView("day")} title="Jour">
      <CalendarSingle size={16}/>
    </button>
    <button class="activityIcon" class:active={view === "agenda"} onclick={() => goToView("agenda")} title="Agenda">
      <List size={16}/>
    </button>
  </nav>

  <aside class="sidebar">
    {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
      <SmallCalendar date={date} smaller={true} onDayClick={(clickedDate) => smallCalendarClick(clickedDate)}></SmallCalendar>
    {/if}

    <div class="sources">
      {@render sourceEntries(repository.sources)}
    </div>

    <Horizontal position="center">
      <IconButton click={showSettingsModal}>
        <Settings/>
      </IconButton>
      {#if !publicReadonly}
        <IconButton click={showSourceWizardModal}>
          <PlusIcon/>
        </IconButton>
        <IconButton click={showCreditsModal}>
          <Copyleft/>
        </IconButton>
      {/if}
    </Horizontal>

    <span class="copyright">
      Copyright © 2026 Kacper Darowski (Opisek)<br>
      Licensed under TBD 
    </span>
  </aside>
</div>

<main>
  <div class="toprow">
    <MonthSelection bind:date granularity={calendarGranularity} />
    <Horizontal position="justify" width="auto">
      {#if connectivity.reachable != Reachability.Database}
        <span class="reachability">
          {#if connectivity.reachable == Reachability.Backend}
            The database cannot be reached.
          {:else if connectivity.reachable == Reachability.Frontend}
            The backend server cannot be reached.
          {:else if connectivity.reachable == Reachability.None}
            The frontend server cannot be reached.
          {:else if connectivity.reachable == Reachability.Incompatible}
            The frontend server and the backend server are not compatible.
          {:else}
            Unknown network error
          {/if}
          <WifiOff size={20}/>
        </span>
      {/if}

      <IconButton click={forceRefresh}>
        <span class="refreshButtonWrapper" class:spin={loaderAnimation} onanimationiteration={() => { if (!isLoading) loaderAnimation = false; }}>
          <RefreshCw size={20}/>
        </span>
      </IconButton>

      {#if !todayInRange}
        <Button onClick={seeToday} compact={true}>
          Today
        </Button>
      {/if}

      {#if !settings.userSettings[UserSettingKeys.ThemeSynchronize]}
        <ThemeToggle/>
      {/if}

    </Horizontal>
  </div>
  {#if view === "agenda"}
    <section class="agenda" bind:this={agendaContainer}>
      {#if agendaDays.length === 0}
        <div class="agendaEmpty">
          Aucun rendez-vous sur cette période.
        </div>
      {:else}
        {#each agendaDays as day (day.date.getTime())}
          <div
            class="agendaDay"
            data-agenda-today={getLocalDayKey(day.date) === todayKey ? "true" : "false"}
          >
            <h3 class="agendaDate">{day.date.toLocaleDateString([], { weekday: "long", day: "2-digit", month: "long", year: "numeric" })}</h3>
            {#each day.events as event (event.id + event.date.start.getTime())}
              {@const agendaOwners = agendaMergedOwnersLabel(event)}
              <div
                class="agendaItem"
                role="button"
                tabindex="0"
                onclick={() => openAgendaEvent(event)}
                onkeydown={(e) => agendaItemKeydown(e, event)}
              >
                <span class="agendaTime">
                  {#if event.date.allDay}
                    Journee
                  {:else}
                    {event.date.start.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                  {/if}
                </span>
                <span class="agendaName">
                  {#each agendaParticipantColors(event) as color (color)}
                    <span class="eventDot" style="background-color: {color}"></span>
                  {/each}
                  <span class="agendaMainLine">
                    <span class="agendaTitle">{event.name}</span>
                    {#if agendaOwners}
                      <span class="agendaOwners" title={agendaOwners}> · {agendaOwners}</span>
                    {/if}
                  </span>
                </span>
              </div>
            {/each}
          </div>
        {/each}
      {/if}
    </section>
  {:else}
    <Calendar
      date={date}
      {today}
      view={view}
      events={repository.events}
      readOnly={publicReadonly}
    />
  {/if}
  <div class="statusbar">
    <span>{timezone}</span>
    <span>{connectivityLabel}</span>
    <span>{nextEventLabel}</span>
  </div>
</main>

{#snippet sourceEntries(sources: SourceModel[])}
  {#each sources as source, i (source.id)}
    <SourceEntry bind:source={repository.sources[i]} readOnly={publicReadonly}/>
    {#if !metadata.collapsedSources.has(repository.sources[i].id)}
      {@render calendarEntries(repository.calendars.filter(cal => cal.source === source.id) || [])}
    {/if}
  {/each}
{/snippet}

{#snippet calendarEntries(calendars: CalendarModel[])}
  {#each calendars as cal (cal.id)}
    {@const index = repository.calendars.findIndex((calendar) => calendar.id === cal.id)}
    <CalendarEntry bind:calendar={repository.calendars[index]} readOnly={publicReadonly}/>
  {/each}
{/snippet}