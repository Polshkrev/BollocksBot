#define PROGRAMME_IMPLEMENTATION
#include "Kada/programme.h"

#define FILES_IMPLEMENTATION
#include "Kada/lib/c/collections/files.h"

#define ENTRY_IMPLEMENTATION
#include "Kada/lib/c/collections/entry.h"

#define FLAG_IMPLEMENTATION
#include "Kada/lib/c/flag.h"

void add_warnings(command_t *command)
{
    command_append(command, " -Wall");
    command_append(command, " -Wextra");
}

void add_compiler_flags(command_t *command, const char *compiler, bool release)
{
    command_append(command, compiler);
    if (!release)
    {
        add_warnings(command);
        command_append(command, " -ggdb");
    }
    command_append(command, " -std=c17 -pedantic");
    if (release) command_append(command, " -O3");
    command_append(command, " -o ");
}

bool add_entry(command_t *command, const entry_t *entry, const char *target_folder, const char *root, const char **result)
{
    if (entry->type != FILE_TYPE) return false;
    path_t parent = path_get_parent(&entry->path);
    if (strcmp(passtr(&parent), root) != 0) return false;
    const char *extension = path_extension(&entry->path);
    if (NULL == extension) return false;
    path_t filename = path_filename(&entry->path);
    const char *filename_string = passtr(&filename);
    if (strcmp(filename_string, root) == 0) return false;
    command_appendf(command, "%s/%s.exe %s", target_folder, filename_string, passtr(&entry->path));
    *result = filename_string;
    return true;
}

int main(int argc, char **argv)
{
    const char *name_flag = "name";
    const char *build_folder = "build";
    const char *bin_folder = "bin";
    const char *compiler = "gcc";
    const char *release_flag = "release";
    char **target_name = flag_string("name", "main", "Set the name of the programme to build.");
    bool *release = flag_bool("release", false, "Set the release mode optimizations.");
    flag_parse(argc, argv);
    logger_t *logger = logger_init(build_folder, LOG_DEBUG);
    logger_add_console(logger);
    programme_t programme = programme_init_with_logger(logger);
    files_t files = files_init(build_folder);
    if (!files_fill(&files))
    {
        logger_log(logger, "RuntimeError: Can not walk directory 'build'.\n", LOG_CRITICAL);
        files_delete(&files);
        programme_delete(&programme);
        return 1;
    }
    for (size_t i = 0; i < files.size; i++)
    {
        command_t build = command_init();
        add_compiler_flags(&build, compiler, *release);
        const char *filename = NULL;
        entry_t entry = entry_init(path_from(*files_at(&files, i)));
        if (!add_entry(&build, &entry, bin_folder, build_folder, &filename)) continue;
        command_t run = command_init();
        command_appendf(&run, "%s\\%s.exe -%s %s", bin_folder, filename, name_flag, *target_name);
        if (*release) command_appendf(&run, " -%s", release_flag);
        programme_append(&programme, &build);
        programme_append(&programme, &run);
    }
    if (!programme_run(&programme))
    {
        logger_log(logger, "RuntimeError: Can not run programme.\n", LOG_CRITICAL);
        files_delete(&files);
        programme_delete(&programme);
        return 1;
    }
    programme_delete(&programme);
    files_delete(&files);
    return 0;
}