<?php

/**
 * @file
 * An example of how to use PHP 7.x+ to communicate with the socket opened by
 * the 'ptpip' command.
 *
 * Execute it like this:
 * ```bash
 * php findAllProperties.php --help
 * ```
 * to see its usage.
 * You can use this command to enumerate all device property codes for your
 * device to see which ones are supported.
 */

global $options, $host, $port, $start, $end, $dryRun;

const OPTIONS = 'h::p::s::e::';
const LONG_OPTS = ['dry-run', 'help'];

const MIN_VALUE = 0;
const MAX_VALUE = 65535;
const DEFAULT_HOST = '127.0.0.1';
const DEFAULT_PORT = 15740;

parseOptions();

for ($i = $start; $i < $end; $i++) {
    $propertyCode = dechex($i);
    $message = sprintf('opreq 0x1015 0x%s', $propertyCode);
    $output = sprintf("Composed message: '%s'", $message);
    $len = strlen($output);
    if (0 === $i) {
        printSeparator($len);
    }
    echo $output . PHP_EOL;
    sendSocketMessageAndPrintResult($host, $port, $message);
    printSeparator($len);
}

function parseOptions(): void
{
    global $options, $host, $port, $start, $end, $dryRun;

    $options = getopt(OPTIONS, LONG_OPTS);

    $help = getOption('help', false, true);
    if ($help) {
        printUsage();
        exit(0);
    }

    $dryRun = getOption('dry-run', false, true);

    $host = getOption('h', DEFAULT_HOST);
    $port = (int) getOption('p', DEFAULT_PORT);

    $start = (int) getOption('s', MIN_VALUE);
    $end = (int) getOption('e', MAX_VALUE);
    if ($port > MAX_VALUE || $start > MAX_VALUE || $end > MAX_VALUE) {
        printf('The maximum allowed value is %d!%s', MAX_VALUE, PHP_EOL);
        exit(1);
    }

    if ($port < MIN_VALUE || $start < MIN_VALUE || $end < MIN_VALUE) {
        printf('The minimum allowed value is %d!%s', MIN_VALUE, PHP_EOL);
        exit(1);
    }
}

function printUsage(): void
{
    $opts = explode('::', OPTIONS);
    foreach ($opts as $opt) {
        if (empty($opt)) {
            continue;
        }

        $message = '';
        switch ($opt) {
            case 'h':
                $message = sprintf('The host to connect to; defaults to %s.', DEFAULT_HOST);
                break;
            case 'p':
                $message = sprintf('The port to connect to; defaults to %d.', DEFAULT_PORT);
                break;
            case 's':
                $message = sprintf('The decimal value to start eumeration from (defaults to %d).', MIN_VALUE);
                break;
            case 'e':
                $message = sprintf('The decimal value to enumerate to (defaults to %d).', MAX_VALUE);
                break;
        }
        echo "-${opt}\n\t$message" . PHP_EOL;
    }

    foreach (LONG_OPTS as $opt) {
        $message = '';
        switch ($opt) {
            case 'dry-run':
                $message = 'Do not actually execute the commands.';
                break;
            case 'help':
                $message = 'Print this message.';
                break;
        }
        echo "--${opt}\n\t$message" . PHP_EOL;
    }
}

/**
 * @param mixed $default The value to return when the option is not set.
 *
 * @return mixed
 */
function getOption(string $option, $default = null, $longOpt = false)
{
    global $options;

    if (array_key_exists($option, $options)) {
        if ($longOpt) {
            return true;
        }
        return $options[$option];
    }

    return $default;
}

function sendSocketMessageAndPrintResult(string $host, int $port, string $message): void
{
    global $dryRun;

    if ($dryRun) {
        echo 'Dry run active, not sending message!' . PHP_EOL;

        return;
    }

    $fp = @fsockopen($host, $port, $errno, $errstr, 30);
    if (!$fp) {
        printf("%s (%d)%s", $errstr, $errno, PHP_EOL);
        exit(1);
    }

    fwrite($fp, $message . "\n");
    while (!feof($fp)) {
        echo fgets($fp, 128);
    }
    fclose($fp);
}

function printSeparator(int $length): void
{
    echo str_repeat('-', $length) . PHP_EOL;
}
