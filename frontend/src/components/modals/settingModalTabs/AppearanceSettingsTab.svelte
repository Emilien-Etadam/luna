<script lang="ts">
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../../types/settings";
  import ColorInput from "../../forms/ColorInput.svelte";
  import SelectInput from "../../forms/SelectInput.svelte";
  import SliderInput from "../../forms/SliderInput.svelte";
  import ToggleInput from "../../forms/ToggleInput.svelte";
  import SectionDivider from "../../layout/SectionDivider.svelte";

  interface Props {
    settings: Settings;
    accentColor: string;
  }

  let {
    settings,
    accentColor = $bindable(),
  }: Props = $props();
</script>

<SectionDivider title={"Calendar Appearance"}/>
<ToggleInput
  name={UserSettingKeys.DisplayAllDayEventsFilled}
  description="Fill All-Day Events"
  bind:value={settings.userSettings[UserSettingKeys.DisplayAllDayEventsFilled]}
/>
<ToggleInput
  name={UserSettingKeys.DisplayNonAllDayEventsFilled}
  description="Fill Non-All-Day Events"
  bind:value={settings.userSettings[UserSettingKeys.DisplayNonAllDayEventsFilled]}
/>
<ToggleInput
  name={UserSettingKeys.DisplaySmallCalendar}
  description="Display Small Calendar"
  bind:value={settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
/>
<ToggleInput
  name={UserSettingKeys.DynamicCalendarRows}
  description="Dynamic Calendar Row Count"
  bind:value={settings.userSettings[UserSettingKeys.DynamicCalendarRows]}
/>
{#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
  <ToggleInput
    name={UserSettingKeys.DynamicSmallCalendarRows}
    description="Dynamic Small Calendar Row Count"
    bind:value={settings.userSettings[UserSettingKeys.DynamicSmallCalendarRows]}
  />
{/if}
<ToggleInput
  name={UserSettingKeys.DisplayWeekNumbers}
  description="Display Week Numbers"
  bind:value={settings.userSettings[UserSettingKeys.DisplayWeekNumbers]}
/>
<SelectInput
  name={UserSettingKeys.FirstDayOfWeek}
  placeholder="First Day of Week"
  bind:value={settings.userSettings[UserSettingKeys.FirstDayOfWeek]}
  options={[
    { name: "Monday", value: 1 },
    { name: "Tuesday", value: 2 },
    { name: "Wednesday", value: 3 },
    { name: "Thursday", value: 4 },
    { name: "Friday", value: 5 },
    { name: "Saturday", value: 6 },
    { name: "Sunday", value: 0 }
  ]}
/>
<SectionDivider title={"Site Appearance"}/>
<ToggleInput
  name={UserSettingKeys.AppearenceFrostedGlass}
  description="Frosted Glass Effect"
  bind:value={settings.userSettings[UserSettingKeys.AppearenceFrostedGlass]}
/>
<ToggleInput
  name={UserSettingKeys.DisplayRoundedCorners}
  description="Rounded Corners"
  bind:value={settings.userSettings[UserSettingKeys.DisplayRoundedCorners]}
/>
<SliderInput
  name={UserSettingKeys.UiScaling}
  title="Scaling"
  bind:value={settings.userSettings[UserSettingKeys.UiScaling]}
  min={0.5}
  max={1.5}
  step={0.05}
  detentTransform={(value) => `${Math.round(value * 100)}%`}
/>
<ColorInput
  name="Accent Color"
  bind:color={accentColor}
  editable={true}
/>
<ToggleInput
  name={UserSettingKeys.ThemeSynchronize}
  description="Synchronize Theme with System"
  bind:value={settings.userSettings[UserSettingKeys.ThemeSynchronize]}
/>
<SectionDivider title="Animations"/>
<SliderInput
  name={UserSettingKeys.AnimationDuration}
  title="Animation Duration"
  info={"To disable animations, set the animation duration to 0%."}
  bind:value={settings.userSettings[UserSettingKeys.AnimationDuration]}
  min={0}
  max={2}
  step={0.1}
  detentTransform={(value) => `${Math.round(value * 100)}%`}
/>
{#if settings.userSettings[UserSettingKeys.AnimationDuration] != 0}
  <ToggleInput
    name={UserSettingKeys.AnimateCalendarSwipe}
    description="Animate Calendar"
    bind:value={settings.userSettings[UserSettingKeys.AnimateCalendarSwipe]}
  />
  {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
    <ToggleInput
      name={UserSettingKeys.AnimateSmallCalendarSwipe}
      description="Animate Small Calendar"
      bind:value={settings.userSettings[UserSettingKeys.AnimateSmallCalendarSwipe]}
    />
  {/if}
  <ToggleInput
    name={UserSettingKeys.AnimateMonthSelectionSwipe}
    description="Animate Month Selection"
    bind:value={settings.userSettings[UserSettingKeys.AnimateMonthSelectionSwipe]}
  />
{/if}