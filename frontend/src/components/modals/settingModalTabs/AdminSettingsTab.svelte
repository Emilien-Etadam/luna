<script lang="ts">
  import { browser } from "$app/environment";
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import { fetchJson, fetchResponse } from "../../../lib/client/net";
  import { queueNotification } from "../../../lib/client/notifications";
  import { ColorKeys } from "../../../types/colors";
  import { GlobalSettingKeys } from "../../../types/settings";
  import SelectButtons from "../../forms/SelectButtons.svelte";
  import ToggleInput from "../../forms/ToggleInput.svelte";
  import Button from "../../interactive/Button.svelte";
  import Paragraph from "../../layout/Paragraph.svelte";
  import SectionTitle from "../../layout/SectionTitle.svelte";

  interface Props {
    settings: Settings;
    showConfirmation: (message: string, onConfirm: () => Promise<void>, confirmText?: string, onCancel?: () => Promise<void>, cancelText?: string) => void;
    refetchProfilePicture: () => void;
    snapshotSettings: () => void;
  }

  let {
    settings,
    showConfirmation,
    refetchProfilePicture,
    snapshotSettings,
  }: Props = $props();

  function resetGlobalSettings() {
    showConfirmation("Are you sure you want to reset all global settings?\nThis action is irreversible.", async () => {
      await fetchResponse("/api/settings", { method: "DELETE" });
      settings.fetchSettings().then(() => {
        snapshotSettings();
        refetchProfilePicture();
      });
    });
  }

  let publicCalEnabled = $state(false);
  let publicCalLoaded = $state(false);

  $effect(() => {
    if (!browser || !settings.userData.admin) {
      return;
    }
    fetchJson<{ enabled: boolean }>("/api/admin/public-calendar")
      .then((d) => {
        publicCalEnabled = d.enabled;
        publicCalLoaded = true;
      })
      .catch(() => {
        publicCalLoaded = true;
      });
  });

  async function onPublicCalendarToggle(next: boolean) {
    try {
      await fetchResponse("/api/admin/public-calendar", {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ enabled: next }),
      });
    } catch {
      publicCalEnabled = !next;
      queueNotification(ColorKeys.Danger, "Impossible de mettre à jour la publication publique.");
    }
  }

  let publicCalendarUrl = $derived(
    typeof window !== "undefined" ? `${window.location.origin}/public` : "",
  );

  async function copyPublicCalendarUrl() {
    const url = typeof window !== "undefined" ? `${window.location.origin}/public` : "";
    if (!url || !navigator.clipboard?.writeText) return;
    await navigator.clipboard.writeText(url);
  }
</script>

<SectionTitle title="Publication publique" />
<Paragraph>
  Activer cette option rend TOUS les événements de TOUS les utilisateurs visibles sans authentification à l’URL publique ci-dessous.
</Paragraph>
<ToggleInput
  name="public_calendar_feed"
  description="Activer la publication publique du calendrier global"
  bind:value={publicCalEnabled}
  editable={publicCalLoaded}
  onChange={onPublicCalendarToggle}
/>
{#if publicCalLoaded && publicCalEnabled && publicCalendarUrl}
  <Paragraph>
    URL publique : <strong>{publicCalendarUrl}</strong>
    <Button color={ColorKeys.Primary} onClick={copyPublicCalendarUrl}>Copier</Button>
  </Paragraph>
{/if}

<ToggleInput
  name={GlobalSettingKeys.RegistrationEnabled}
  description="Enable Open Registration"
  info={"Allows anyone to create an account.\nIf you just want to invite a few people, head to the \"Users\" tab."}
  bind:value={settings.globalSettings[GlobalSettingKeys.RegistrationEnabled]}
/>
<!--
<ToggleInput
  name={GlobalSettingKeys.UseCdnFonts}
  description="Use Google's CDN for fonts"
  bind:value={settings.globalSettings[GlobalSettingKeys.UseCdnFonts]}
/>
-->
<ToggleInput
  name={GlobalSettingKeys.UseIpGeolocation}
  description="Determine IP Address Geolocation"
  info={`This setting enables the frontend to display an approximate location of the user's sessions in their settings page. This feature is meant to help users determine whether illicit actors have access to their account. Since determining the geolocation of an IP address requires the use of a third-party service, using it raises some privacy concerns. If you would prefer not to use this feature, it can be disabled here. Please note that IP addresses themselves will still be stored in the database for security and audit purposes, but they are never shared with any third party outside of the behavior controlled by this toggle.`}
  bind:value={settings.globalSettings[GlobalSettingKeys.UseIpGeolocation]}
/>
<ToggleInput
  name={GlobalSettingKeys.EnableGravatar}
  description="Enable Gravatar Profile Pictures"
  info={`Gravatar is a third party service that lets users set their profile picture once and have it associated with their e-mail address.  This is accomplished by querying said server with an MD5 hash of the user's e-mail address. This can be seen as a privacy concern, as it allows Gravatar to track users across different services. If you would prefer not to use Gravatar profile pictures, you can disable them here. This setting disables the use of Gravatar site-wide.`}
  bind:value={settings.globalSettings[GlobalSettingKeys.EnableGravatar]}
/>
<ToggleInput
  name={GlobalSettingKeys.CacheProfilePictures}
  description="Cache Profile Pictures"
  info={`Whether profile pictures from remote websites (including Gravatar) should be cached locally. Caching profile pictures can improve performance and privacy, but it introduces the risk of hosting illicit content if a malicious user uses it as their profile picture. The same risk exists if users upload profile pictures directly.`}
  bind:value={settings.globalSettings[GlobalSettingKeys.CacheProfilePictures]}
/>
<ToggleInput
  name={GlobalSettingKeys.EnableProfilePicturesUpload}
  description="Enable Profile Picture Uploads"
  info={`Allow users to upload their own profile pictures directly to the server. This introduces the risk of hosting illicit content if a malicious user uploads such an image as their profile picture. The same risk exists if caching of remote profile pictures is enabled.`}
  bind:value={settings.globalSettings[GlobalSettingKeys.EnableProfilePicturesUpload]}
/>
<SelectButtons
  name={GlobalSettingKeys.LoggingVerbosity}
  bind:value={settings.globalSettings[GlobalSettingKeys.LoggingVerbosity]}
  placeholder="Error Messages Verbosity"
  info={"How much information about errors is returned to the user.\nThis setting applies to all users.\n\"Debug\" should never be used in production."}
  options={[
    { name: "Broad", value: 3 },
    { name: "Plain", value: 2 },
    { name: "Wordy", value: 1 },
    { name: "Debug", value: 0 }
  ]}
/>
<Button color={ColorKeys.Danger} onClick={resetGlobalSettings}>Reset all global settings</Button>