import { ConfigurableModuleBuilder } from '@nestjs/common';
import { ModuleOptions } from './interfaces/module-options';

export const { ConfigurableModuleClass, MODULE_OPTIONS_TOKEN } =
  new ConfigurableModuleBuilder<ModuleOptions>().build();
