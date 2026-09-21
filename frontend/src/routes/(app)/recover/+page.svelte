<script lang="ts">
  import Link from '../../../components/forms/Link.svelte';
  import SimplePage from '../../../components/layout/SimplePage.svelte';
  import Title from '../../../components/layout/Title.svelte';
  import Paragraph from '../../../components/layout/Paragraph.svelte';
  import { page } from '$app/state';

  const redirect = $derived(page.url.searchParams.get('redirect') || "/");
</script>

<style lang="scss">
  @use "../../../styles/dimensions.scss";
  @use "../../../styles/text.scss";

  .panel {
    border-radius: var(--borderRadiusLarge);
    max-width: 50vw;
    min-width: 30em;
    padding: 24px 28px 28px 28px;
    background-color: var(--surface-raised);
    border: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapMiddle;
    box-shadow: var(--shadow-2);
  }

  @media (max-width: 720px) {
    .panel {
      max-width: 100vw;
      min-width: 0;
      width: 100%;
      padding: 20px;
    }
  }

  pre {
    margin: 0;
    padding: dimensions.$gapMiddle;
    overflow-x: auto;
    font-size: text.$fontSizeSmall;
    line-height: text.$lineHeightParagraph;
    white-space: pre;
    background-color: var(--surface-sunken);
    border: 1px solid var(--border-default);
    border-radius: var(--borderRadius);
  }
</style>

<SimplePage>
  <div class="panel">
    <Title>Forgot password</Title>
    <Paragraph>
      Luna cannot recover the original password: it is stored as an irreversible hash.
      Only the instance operator can set a new one from the server.
    </Paragraph>
    <Paragraph>
      If you are a regular user, contact the administrator of this instance.
    </Paragraph>
    <Paragraph>
      If you operate this instance, list accounts to find the admin username, then reset the password.
      Docker:
    </Paragraph>
    <pre>docker exec -it luna-backend ./luna-backend users
docker exec -it luna-backend ./luna-backend reset-password YOUR_USERNAME</pre>
    <Paragraph>
      Bare metal, from backend/src with the same environment as the service:
    </Paragraph>
    <pre>./luna-backend users
./luna-backend reset-password YOUR_USERNAME</pre>
    <Paragraph>
      Without a TTY, pass the new password via LUNA_NEW_PASSWORD. If no account is administrator:
    </Paragraph>
    <pre>luna-backend promote-admin YOUR_USERNAME</pre>
    <Link href="/login?redirect={encodeURIComponent(redirect)}">Back to login</Link>
  </div>
</SimplePage>
