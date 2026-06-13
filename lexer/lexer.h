#ifndef LEXER_H_
#define LEXER_H_

#define COMMAND_IMPLEMENTATION
#include "bot_command.h"

/**
 * @brief A lexer.
 */
typedef struct
{
    const char *source;
    size_t cursor;
} lexer_t;

/**
 * @brief Construct a new lexer.
 * @returns A new lexer.
 */
lexer_t lexer_init(void);

/**
 * @brief Set the source of the lexer.
 * @param lexer Lexer to which to set the source.
 * @param source Source from which to set within the lexer.
 */
void lexer_set_source(lexer_t *lexer, const char *source);

/**
 * @brief Peek at the current state of the lexer.
 * @param lexer Lexer from which to peek.
 * @returns The character stored within the lexer at its cursor.
 */
char lexer_peek(const lexer_t *lexer);

/**
 * @brief Advance the cursor of the lexer by a factor of one.
 * @param lexer Lexer to increment.
 */
void lexer_advance(lexer_t *lexer);

/**
 * @brief Consume all whitespace characters at the current cursor.
 * @param lexer Lexer to which to consume the whitespace.
 */
void lexer_skip_whitespace(lexer_t *lexer);

/**
 * @brief Obtain the length of the word at the current cursor position.
 * @param lexer Lexer from which to read the length of the word.
 * @returns The length of the word at the cursors current position.
 */
size_t lexer_get_word_length(lexer_t *lexer);

/**
 * @brief Obtain the algorythmic next token.
 * @param lexer Lexer from which to read the next token.
 * @returns The next significant token within the lexer's source.
 */
token_t lexer_next_token(lexer_t *lexer);

/**
 * @brief Represent the given lexer's source as an array of tokens.
 * @param lexer Lexer from which to tokenize.
 * @param output Out parametre which initializes an array of tokens.
 */
void tokenize(lexer_t *lexer, token_array_t *output);

#endif // LEXER_H_

#ifdef LEXER_IMPLEMENTATION

#include <ctype.h> // isspace

/**
 * @brief Construct a new lexer.
 * @returns A new lexer.
 */
lexer_t lexer_init(void)
{
    return (lexer_t)
    {
        .cursor = 0
    };
}

/**
 * @brief Set the source of the lexer.
 * @param lexer Lexer to which to set the source.
 * @param source Source from which to set within the lexer.
 */
void lexer_set_source(lexer_t *lexer, const char *source)
{
    lexer->source = source;
}

/**
 * @brief Peek at the current state of the lexer.
 * @param lexer Lexer from which to peek.
 * @returns The character stored within the lexer at its cursor.
 */
char lexer_peek(const lexer_t *lexer)
{
    return lexer->source[lexer->cursor];
}

/**
 * @brief Advance the cursor of the lexer by a factor of one.
 * @param lexer Lexer to increment.
 */
void lexer_advance(lexer_t *lexer)
{
    lexer->cursor++;
}

/**
 * @brief Consume all whitespace characters at the current cursor.
 * @param lexer Lexer to which to consume the whitespace.
 */
void lexer_skip_whitespace(lexer_t *lexer)
{
    while (isspace(lexer_peek(lexer)))
    {
        lexer_advance(lexer);
    }
}

/**
 * @brief Obtain the length of the word at the current cursor position.
 * @param lexer Lexer from which to read the length of the word.
 * @returns The length of the word at the cursors current position.
 */
size_t lexer_get_word_length(lexer_t *lexer)
{
    size_t length = 0;
    while (lexer->source[lexer->cursor + length] && !isspace(lexer->source[lexer->cursor + length]) && lexer->source[lexer->cursor + length] != *token_kind_to_string(TokenColon))
    {
        length++;
    }
    return length;
}

/**
 * @brief Obtain the algorythmic next token.
 * @param lexer Lexer from which to read the next token.
 * @returns The next significant token within the lexer's source.
 */
token_t lexer_next_token(lexer_t *lexer)
{
    lexer_skip_whitespace(lexer);
    char character = lexer_peek(lexer);
    const char *start = &character;

    if (!character) return token_init(TokenEnd, start, 0);
    else if (character == *token_kind_to_string(TokenColon))
    {
        lexer_advance(lexer);                                               
        return token_init(TokenColon, start, 1);
    }
    const char *word = &lexer->source[lexer->cursor];
    size_t length = lexer_get_word_length(lexer);

    if (length > MAX_TOKENS) return token_init(TokenEnd, word, length);
    
    lexer->cursor += length;
    return token_init(TokenWord, word, length);
}

/**
 * @brief Represent the given lexer's source as an array of tokens.
 * @param lexer Lexer from which to tokenize.
 * @param output Out parametre which initializes an array of tokens.
 */
void tokenize(lexer_t *lexer, token_array_t *output)
{
    token_t token = {0};
    do
    {
        if (output->size >= MAX_TOKENS) break;
        token = lexer_next_token(lexer);
        output->tokens[output->size++] = token;
    }
    while (token.kind != TokenEnd);
}

#endif // LEXER_IMPLEMENTATION