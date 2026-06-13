#define PROGRAMME_IMPLEMENTATION
#include "Kada/programme.h"

#define FLAG_IMPLEMENTATION
#include "Kada/lib/c/flag.h"

void add_linker_flags(command_t *command)
{
    command_append(command, " -ldflags=\"-s -w\"");
}

int main(int argc, char **argv)
{
    char **target_name = flag_string("name", "main", "Set the name of the programme to compile.");
    flag_parse(argc, argv);
    const char *compiler = "go";
    const char *name = "compile";
    const char *command_folder = "cmd";
    const char *target_folder = "bin";
    logger_t *logger = logger_init(name, LOG_DEBUG);
    logger_add_console(logger);
    command_t compile = command_init();
    command_append(&compile, compiler);
    command_append(&compile, " build");
    add_linker_flags(&compile);
    command_appendf(&compile, " -o %s/%s.exe %s/%s.go", target_folder, *target_name, command_folder, *target_name);
    programme_t programme = programme_init_with_logger(logger);
    programme_append(&programme, &compile);
    if (!programme_run(&programme))
    {
        size_t checkpoint = buffer_save();
        logger_log(logger, buffer_sprintf("RuntimeError: Can not run command: '%s'.", command_data(&compile)), LOG_CRITICAL);
        logger_delete(logger);
        command_delete(&compile);
        buffer_rewind(checkpoint);
        return 1;
    }
    logger_delete(logger);
    command_delete(&compile);
    return 0;
}