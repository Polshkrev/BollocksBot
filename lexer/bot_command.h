#define MESSAGE_IMPLEMENTATION
#include "message.h" // message_t

/**
 * @brief Detemine if the given command is valid.
 * @param command Command which to analyse.
 * @returns True if the given command is a valid irc command.
 */
bool command_is_valid(const string_t *command);

/**
 * @brief Extract the name of the given message.
 * @param message Message from which to parse.
 * @param index Index of the parse message.
 * @returns A sized string containing the name of the command.
 */
string_t command_extract_name(const string_t *message, size_t *index);

/**
 * @brief Extract the argument of the given message.
 * @param message Message from which to parse.
 * @param index Index of the parse message.
 * @returns A sized string containing the argument of the command.
 */
string_t command_extract_arguments(const string_t *message, size_t *index);

/**
 * @brief Parse a given command.
 * @param message Command to parse.
 * @returns True if the command is able to be parsed, else false.
 */
bool parse_command(message_t *message);

#ifdef COMMAND_IMPLEMENTATION

#include <ctype.h> // isspace

/**
 * @brief Detemine if the given command is valid.
 * @param command Command which to analyse.
 * @returns True if the given command is a valid irc command.
 */
bool command_is_valid(const string_t *command)
{
    if (!command->count || !command->data) return false;
    else if (command->count < 1) return false;
    else if (command->data[0] != *token_type_to_string(TokenBang)) return false;
    return true;
}

/**
 * @brief Extract the name of the given message.
 * @param message Message from which to parse.
 * @param index Index of the parse message.
 * @returns A sized string containing the name of the command.
 */
string_t command_extract_name(const string_t *message, size_t *index)
{
    size_t start = *index + 1;
    while (message->data[*index] && !isspace(message->data[*index])) (*index)++;
    return string_new(message->data + start, *index - start);
}

/**
 * @brief Extract the argument of the given message.
 * @param message Message from which to parse.
 * @param index Index of the parse message.
 * @returns A sized string containing the argument of the command.
 */
string_t command_extract_arguments(const string_t *message, size_t *index)
{
    while (*index < message->count && isspace(message->data[*index])) (*index)++;
    if (*index >= message->count) return string_null;
    if (message->data[*index] != *token_type_to_string(TokenDoubleQuote) && message->data[*index] != *token_type_to_string(TokenSingleQuote)) return string_new(message->data + *index, message->count - *index);
    (*index)++;
    size_t start = *index;
    while (*index < message->count && message->data[*index] != *token_kind_to_string(TokenDoubleQuote) && message->data[*index] != *token_kind_to_string(TokenSingleQuote)) (*index)++;
    return string_new(message->data + start, *index - start);
}

/**
 * @brief Parse a given command.
 * @param message Command to parse.
 * @returns True if the command is able to be parsed, else false.
 */
bool parse_command(message_t *message)
{
    if (!command_is_valid(&message->text)) return false;
    size_t index = 0;
    message->command = command_extract_name(&message->text, &index);
    message->arguments = command_extract_arguments(&message->text, &index);
    return true;
}

#endif // COMMAND_IMPLEMENTATION