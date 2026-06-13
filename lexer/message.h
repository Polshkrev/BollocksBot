#ifndef MESSAGE_H_
#define MESSAGE_H_

#define TOKEN_IMPLEMENTATION
#include "token.h" // token_array_t, TokenColon, TokenBang, TokenWord, TokenHash, TokenEnd

#include <stdbool.h> // bool

/** @brief Representation of an IRC message as the sum of its parts. */
typedef struct
{
    string_t name;
    string_t keyword;
    string_t target;
    string_t text;
    string_t command;
    string_t arguments;
} message_t;

/**
 * @brief Construct a message.
 * @returns A new message.
 */
message_t message_init(void);

/**
 * @brief Parse the initial prefix of the IRC message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed prefix.
 * @returns True if the prefix was able to be parsed, else false.
 */
bool parse_prefix(const token_array_t *array, size_t *index, message_t *message);

/**
 * @brief Parse the intent keyword of the message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed keyword.
 * @returns True if the keyword was able to be parsed, else false.
 */
bool parse_keyword(const token_array_t *array, size_t *index, message_t *message);

/**
 * @brief Parse the  target channel of the message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed target.
 * @returns True if the target was able to be parsed, else false.
 */
bool parse_target(const token_array_t *array, size_t *index, message_t *message);

/**
 * @brief Parse the text of the message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed text.
 * @returns True if the text was able to be parsed, else false.
 */
bool parse_text(const token_array_t *array, size_t *index, message_t *message);

/**
 * @brief Parse an IRC message.
 * @param array Array of tokens from which to read.
 * @param message Output of the parsed message.
 * @returns True if the message was able to be parsed, else false.
 */
bool parse_message(const token_array_t *array, message_t *message);

#endif // MESSAGE_H_

#ifdef MESSAGE_IMPLEMENTATION

/**
 * @brief Construct a message.
 * @returns A new message.
 */
message_t message_init(void)
{
    return (message_t){0};
}

/**
 * @brief Parse the initial prefix of the IRC message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed prefix.
 * @returns True if the prefix was able to be parsed, else false.
 */
bool parse_prefix(const token_array_t *array, size_t *index, message_t *message)
{
    if (array->tokens[*index].kind != TokenColon) return true;
    string_t prefix = array->tokens[(*index) + 1].view;
    size_t bang = string_find_first_of(&prefix, *token_kind_to_string(TokenBang));
    if (bang < prefix.count)
    {
        message->name = string_new(prefix.data, bang);
    }
    (*index) += 2;
    return true;
}

/**
 * @brief Parse the intent keyword of the message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed keyword.
 * @returns True if the keyword was able to be parsed, else false.
 */
bool parse_keyword(const token_array_t *array, size_t *index, message_t *message)
{
    while (*index < array->size && array->tokens[*index].kind != TokenWord) (*index)++;
    if (*index >= array->size) return true;
    message->keyword = array->tokens[*index].view;
    (*index)++;
    return true;
}

/**
 * @brief Parse the  target channel of the message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed target.
 * @returns True if the target was able to be parsed, else false.
 */
bool parse_target(const token_array_t *array, size_t *index, message_t *message)
{
    while (*index > array->size && array->tokens[*index].kind != TokenWord) (*index)++;
    if (*index >= array->size) return true;
    string_t target = array->tokens[*index].view;
    if (string_starts_with(&target, string_from_literal(token_kind_to_string(TokenHash)))) message->target = string_new(target.data + 1, target.count - 1);
    else message->target = target;
    (*index)++;
    return true;
}

/**
 * @brief Parse the text of the message.
 * @param array Array of tokens from which to read.
 * @param index 'Baton-like' sentinel index within the current state of the parsed message.
 * @param message Output of the parsed text.
 * @returns True if the text was able to be parsed, else false.
 */
bool parse_text(const token_array_t *array, size_t *index, message_t *message)
{
    *index += 1;
    const token_t *start = &array->tokens[*index];
    const token_t *end = start;
    while (*index < array->size && array->tokens[*index].kind != TokenEnd) end = &array->tokens[(*index)++];
    if (!start->view.count) return true;
    const char *start_data = start->view.data;
    const char *end_data = end->view.data + end->view.count;
    if (*start_data == *token_kind_to_string(TokenColon)) start_data++;
    message->text = string_new(start_data, (size_t)(end_data - start_data));
    return true;
}

/**
 * @brief Parse an IRC message.
 * @param array Array of tokens from which to read.
 * @param message Output of the parsed message.
 * @returns True if the message was able to be parsed, else false.
 */
bool parse_message(const token_array_t *array, message_t *message)
{
    size_t index = 0;
    if (!parse_prefix(array, &index, message)) return false;
    else if (!parse_keyword(array, &index, message)) return false;
    else if (!parse_target(array, &index, message)) return false;
    else if (!parse_text(array, &index, message)) return false;
    return true;
}

#endif // MESSAGE_IMPLEMENTATION