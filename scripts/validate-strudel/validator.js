import * as readline from 'readline';

let transpiler;

async function init() {
  try {
    const strudelTranspiler = await import('@strudel/transpiler');
    transpiler = strudelTranspiler.transpiler;
  } catch (err) {
    console.error(
      JSON.stringify({
        error: `failed to load @strudel/transpiler: ${err.message}`,
        fatal: true,
      })
    );

    process.exit(1);
  }
}

function validate(code) {
  if (!code || typeof code !== 'string') {
    return { valid: false, error: 'Empty or invalid code' };
  }

  try {
    transpiler(code, {
      wrapAsync: true,
      addReturn: true,
      simpleLocs: true,
    });

    return { valid: true };
  } catch (err) {
    return {
      valid: false,
      error: err.message,
      line: err.loc?.line || err.location?.start?.line || null,
      column: err.loc?.column || err.location?.start?.column || null,
    };
  }
}

async function main() {
  await init();
  console.log(JSON.stringify({ ready: true }));

  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false,
  });

  rl.on('line', line => {
    try {
      const request = JSON.parse(line);

      switch (request.type) {
        case 'validate':
          const result = validate(request.code);
          console.log(JSON.stringify({ id: request.id, ...result }));
          break;
        case 'ping':
          console.log(JSON.stringify({ id: request.id, pong: true }));
          break;
        default:
          console.log(JSON.stringify({ id: request.id, error: 'unknown request type' }));
          break;
      }
    } catch (err) {
      console.log(JSON.stringify({ error: `parse error: ${err.message}` }));
    }
  });

  rl.on('close', () => {
    process.exit(0);
  });
}

main();
