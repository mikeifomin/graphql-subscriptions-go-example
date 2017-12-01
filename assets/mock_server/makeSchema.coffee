fs = require 'fs'
{makeExecutableSchema} = require 'graphql-tools'

resolvers = require './resolvers.coffee'

FILENAME = __dirname + "/../../../schema.go"

START = 'start_used_for_js_mock_server'
END = 'end_used_for_js_mock_server'

module.exports = ()->
  file = fs.readFileSync(FILENAME)
  schema = file.toString().match(/`([^`]+)`/)?[1]
  makeExecutableSchema {
    typeDefs: [schema]
    resolvers
  }
