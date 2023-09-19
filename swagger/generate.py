from swagger.builder import BackboneSwagger


command_swagger = BackboneSwagger(name="bb-swagger", version="v1", generateMap=True)
command_swagger.save_swagger()
command_swagger.save_commandApiMap()

# query_swagger = BackboneSwagger(name="queries", version="v1", config_file="queries.json")
# query_swagger.save_swagger()