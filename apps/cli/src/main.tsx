#!/usr/bin/env node
import React from 'react';
import {render} from 'ink';
import meow from 'meow';
import App from './app.js';
import {readConfig} from '@ordo/core';

readConfig()
const cli = meow(
	`
	Usage
	  $ ordo-cli

	Options
		--name  Your name

	Examples
	  $ ordo-cli --name=Jane
	  Hello, Jane
`,
	{
		importMeta: import.meta,
		flags: {
			name: {
				type: 'string',
			},
		},
	},
);

render(<App name={cli.flags.name} />);
