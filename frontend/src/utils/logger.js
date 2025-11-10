import { info, warn, error, debug, trace, attachConsole } from '@tauri-apps/plugin-log';

/**
 * Logger utility untuk Tauri app
 * Semua log akan tersimpan di folder logs terpisah untuk troubleshooting
 */

// Attach console untuk forward semua console.log ke Tauri log
let detachConsole = null;

/**
 * Initialize logger - attach console untuk forward logs
 */
export async function initLogger() {
  try {
    // Check if running in Tauri
    if (typeof window !== 'undefined' && '__TAURI_INTERNALS__' in window) {
      detachConsole = await attachConsole();
      info('Logger initialized - console logs will be forwarded to log files');
    }
  } catch (err) {
    console.error('Failed to initialize logger:', err);
  }
}

/**
 * Detach console logger
 */
export function detachLogger() {
  if (detachConsole) {
    detachConsole();
    detachConsole = null;
  }
}

/**
 * Log info message
 */
export function logInfo(message, ...args) {
  const logMessage = args.length > 0 ? `${message} ${JSON.stringify(args)}` : message;
  info(logMessage);
  console.log(`[INFO] ${logMessage}`);
}

/**
 * Log warning message
 */
export function logWarn(message, ...args) {
  const logMessage = args.length > 0 ? `${message} ${JSON.stringify(args)}` : message;
  warn(logMessage);
  console.warn(`[WARN] ${logMessage}`);
}

/**
 * Log error message
 */
export function logError(message, error = null, ...args) {
  let logMessage = message;
  
  if (error) {
    logMessage += ` - Error: ${error.message || error}`;
    if (error.stack) {
      logMessage += `\nStack: ${error.stack}`;
    }
  }
  
  if (args.length > 0) {
    logMessage += ` ${JSON.stringify(args)}`;
  }
  
  error(logMessage);
  console.error(`[ERROR] ${logMessage}`);
}

/**
 * Log debug message
 */
export function logDebug(message, ...args) {
  const logMessage = args.length > 0 ? `${message} ${JSON.stringify(args)}` : message;
  debug(logMessage);
  console.debug(`[DEBUG] ${logMessage}`);
}

/**
 * Log trace message
 */
export function logTrace(message, ...args) {
  const logMessage = args.length > 0 ? `${message} ${JSON.stringify(args)}` : message;
  trace(logMessage);
  console.trace(`[TRACE] ${logMessage}`);
}

/**
 * Log API request
 */
export function logApiRequest(method, url, data = null) {
  const message = `API Request: ${method} ${url}`;
  const logData = data ? ` - Data: ${JSON.stringify(data)}` : '';
  info(`${message}${logData}`);
}

/**
 * Log API response
 */
export function logApiResponse(method, url, status, data = null) {
  const message = `API Response: ${method} ${url} - Status: ${status}`;
  const logData = data ? ` - Data: ${JSON.stringify(data)}` : '';
  info(`${message}${logData}`);
}

/**
 * Log API error
 */
export function logApiError(method, url, error) {
  const message = `API Error: ${method} ${url}`;
  logError(message, error);
}

export default {
  init: initLogger,
  detach: detachLogger,
  info: logInfo,
  warn: logWarn,
  error: logError,
  debug: logDebug,
  trace: logTrace,
  apiRequest: logApiRequest,
  apiResponse: logApiResponse,
  apiError: logApiError,
};

