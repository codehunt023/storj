<script lang="ts">
  import type { Operation, ParamUI, InputText, Select } from "./ui-generator";

  import UIInputText from "./UIGeneratorInputText.svelte";
  import UISelect from "./UIGeneratorSelect.svelte";

  function isInputText(p: ParamUI): p is InputText {
    return (p as InputText).type !== undefined;
  }

  function isSelect(p: ParamUI): p is Select {
    return (p as Select).multiple !== undefined;
  }

  export let operation: Operation;
  let params: (boolean | number | string)[] = new Array(
    operation.params.length
  );

  // TODO: the on:click function in the button must not fire if the required
  // fiels aren't set.
</script>

<div>
  <p>{operation.desc}</p>
  <form action="">
    {#each operation.params as param, i}
      <br />
      {#if isInputText(param[1])}
        <UIInputText
          label={param[0]}
          config={param[1]}
          bind:value={params[i]}
        />
      {/if}
      {#if isSelect(param[1])}
        <UISelect label={param[0]} config={param[1]} bind:value={params[i]} />
      {/if}
    {/each}
    <br />
    <button on:click={() => operation.func(...params)}>submit</button>
  </form>
</div>
<pre>
  {#each params as param}
    {param}
  {/each}
</pre>

<style>
</style>
