import { RestScriptEditor } from './restscript/editor';

function Setup(elId: string, updateCallback: (x: string) => void) {
  let codeEl = window.document.getElementById(elId);
  let rse = new RestScriptEditor(codeEl, true);
  rse.updateCallback = updateCallback;
  return rse;
}

eval("window.rseSetup = Setup");

export default Setup;
