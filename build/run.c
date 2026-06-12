#define COMMAND_IMPLEMENTATION
#include "Kada/command.h"

#define FLAG_IMPLEMENTATION
#include "Kada/lib/c/flag.h"

int main(int argc, char **argv)
{
    char **target_name = flag_string("name", "main", "Set the name of the programme to run.");
    flag_parse(argc, argv);
    const char *name = "run";
    const char *target_folder = "bin";
    logger_t *logger = logger_init(name, LOG_DEBUG);
    logger_add_console(logger);
    command_t run = command_init();
    command_appendf(&run, "%s/%s.exe", target_folder, *target_name);
    process_t process = command_run_async_logged(&run, logger);
    if (!process_wait(process))
    {
        size_t checkpoint = buffer_save();
        logger_log(logger, buffer_sprintf("RuntimeError: Can not run command: %s.\n", command_data(&run)), LOG_CRITICAL);
        command_delete(&run);
        logger_delete(logger);
        buffer_rewind(checkpoint);
        return 1;
    }
    logger_delete(logger);
    command_delete(&run);
    return 0;
}