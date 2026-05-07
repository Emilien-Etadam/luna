<script lang="ts">
  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getDayIndex, isSameDay } from "$lib/common/date";
  import { UserSettingKeys } from "../../types/settings";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { setContext } from "svelte";
  import { svelteFlyInHorizontal, svelteFlyOutHorizontal } from "$lib/client/animations";

  interface Props {
    date: Date;
    onDayClick?: (date: Date) => any;
    smaller?: boolean;
  }

  let {
    date = $bindable(new Date()),
    onDayClick = NoOp,
    smaller = false,
  }: Props = $props();

  const settings = getSettings();

  let today = new Date();

  /* Date calculation */
  let [days, amountOfRows] = $derived.by(() => {
    const firstMonthDay = new Date(date.getFullYear(), date.getMonth(), 1);
    const lastMonthDay = new Date(date.getFullYear(), date.getMonth() + 1, 0);
    const firstDayOfWeek = getDayIndex(firstMonthDay);

    const amountOfRows = 
      settings.userSettings[UserSettingKeys.DynamicSmallCalendarRows] ?
      Math.ceil((lastMonthDay.getDate() + firstDayOfWeek) / 7)
      : 6;

    const firstViewDay = new Date(firstMonthDay);
    firstViewDay.setDate(firstMonthDay.getDate() - firstDayOfWeek);
    const lastViewDay = new Date(firstMonthDay);
    lastViewDay.setDate(firstMonthDay.getDate() + 7 * amountOfRows - 1);

    // Fill
    const days = [];

    const dateIterator = new Date(firstViewDay);

    for (let i = 0; i < 7 * amountOfRows; i++) {
      days.push(new Date(dateIterator));
      dateIterator.setDate(dateIterator.getDate() + 1);
    }

    return [days, amountOfRows];
  });

  /* Animation */
  let viewIteration = $state(0);
  // TODO: why do we need displayDays here but not in the large calendar?
  // svelte-ignore state_referenced_locally
  let displayDays = $state(days);
  let currentDate = $state(new Date(date));
  let flyDirection = $state("left");
  setContext("flyDirection", () => flyDirection);
  $effect(() => {
    if (date.getTime() === currentDate.getTime()) return;
    flyDirection = currentDate.getTime() <= date.getTime() ? "left" : "right";
    currentDate = new Date(date);
    viewIteration++;
    displayDays = days;
  });

</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.animation {
    overflow: hidden;
    position: relative;
  }

  div.calendar {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 2px;
    width: 100%;
  }

  div.calendar.animate:not(:first-child) {
    position: absolute;
    top: 0;
    left: 0;
  }

  div.smaller {
    font-size: text.$fontSizeSmall;
    gap: 2px;
  }

  button.day {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: var(--radius-1);
    font-size: var(--font-size-day-number);
    font-variant-numeric: tabular-nums;
    font-weight: 400;
    color: var(--fg-primary);
    background-color: transparent;
    padding: 3px 0;
    cursor: pointer;
    user-select: none;
    position: relative;
    overflow: visible;
    transition: background-color var(--transition-fast), color var(--transition-fast);
  }

  button.day:hover {
    background-color: var(--bg-hover);
  }

  button.day:focus-visible {
    box-shadow: var(--focus-ring);
  }

  button.day.weekend {
    color: var(--fg-muted);
  }

  button.day.today {
    color: var(--fg-inverted);
    font-weight: var(--font-weight-semibold);
    background-color: var(--accent-blue);
    --barFocusIndicatorColor: var(--accent-blue);
  }

  button.day.today:hover {
    background-color: var(--accent-blue);
    filter: brightness(1.08);
  }

  button.day.otherMonth {
    color: var(--fg-disabled);
    opacity: 1;
  }
  button.day.today.otherMonth {
    color: var(--fg-inverted);
  }
</style>

{#if settings.userSettings[UserSettingKeys.AnimateSmallCalendarSwipe]}
  <div class="animation">
    {#each [ displayDays ] as currentDays (viewIteration)}
      {@render grid(currentDays, amountOfRows, true)}
    {/each}
  </div>
{:else}
  {@render grid(days, amountOfRows, false)}
{/if}

{#snippet grid(days: Date[], amountOfRows: number, animate: boolean)}
  <div
    class="calendar"
    class:smaller={smaller}
    class:animate={animate}
    style="grid-template-rows: repeat({amountOfRows}, 1fr)"
    in:svelteFlyInHorizontal={{duration: animate ? 100 : 0}}
    out:svelteFlyOutHorizontal={{duration: animate ? 100 : 0}}
  >
    {#each days as day}
      <button
        class="day"
        class:weekend={day.getDay() === 0 || day.getDay() === 6}
        class:today={isSameDay(day, today)}
        class:otherMonth={day.getMonth() != currentDate.getMonth()}
        type="button"
        onclick={() => (onDayClick(day))}
        use:focusIndicator
      >
        {day.getDate()}
      </button>
    {/each}
  </div>
{/snippet}