package main

import (
		"fmt"
			"go/ast"
				"go/parser"
					"go/token"
						"path/filepath"
					)

					func findImportedStructObjects(filename string) []map[string]string  {
							// Set up the parser configuration
								fset := token.NewFileSet()
									parserConfig := &parser.Config{Mode: parser.ParseComments}
										// Parse the file and extract the import statements
											file, err := parserConfig.ParseFile(fset, filename, nil, parser.ParseComments)
												if err != nil  {
															fmt.Printf("Error parsing file '%s': %v", filename, err)

																	return nil

																		}

																			importedStructObjects := []map[string]string{}

																				for _, imp := range file.Imports  {

																							// Get the package name and file path for the imported package

																									pkgName, pkgPath := filepath.Split(imp.Path.Value[1 : len(imp.Path.Value)-1])

																											// Parse the imported package to find struct objects

																													importsConfig := parserConfig
																															importsConfig.Mode |= parser.ImportsOnly
																																	importsConfig.ParseFile(pkgPath, nil, nil, parser.ParseComments)

																																			pkg, err := parser.ParseDir(fset, pkgPath, nil, parser.ParseComments)

																																					if err != nil  {

																																									fmt.Printf("Error parsing package '%s': %v", pkgPath, err)

																																												continue

																																														}

																																																for _, file := range pkg  {

																																																				for _, decl := range file.Decls  {

																																																									if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE  {

																																																															for _, spec := range genDecl.Specs  {

																																																																						if typeSpec, ok := spec.(*ast.TypeSpec); ok  {

																																																																														if _, ok := typeSpec.Type.(*ast.StructType); ok  {

																																																																																							// Found an imported struct object

																																																																																															structName := typeSpec.Name.Name

																																																																																																							importedStructObjects = append(importedStructObjects, map[string]string {

																																																																																																																	"struct_name": structName,

																																																																																																																										"package_name": pkgName[:len(pkgName)-1],

																																																																																																																																			"file_name": pkgPath,
																																																																																																																																											})

																																																																																																																																																		}

																																																																																																																																																								}

																																																																																																																																																													}

																																																																																																																																																																	}

																																																																																																																																																																				}

																																																																																																																																																																						}

																																																																																																																																																																							}

																																																																																																																																																																								return importedStructObjects

																																																																																																																																																																							}

																																																																																																																																																																							func main()  {

																																																																																																																																																																									// Example usage: find all imported struct objects in a file called "example.go"

																																																																																																																																																																										importedStructObjects := findImportedStructObjects(`C:\\Users\\Gleb\\Desktop\\Учеба\\Диплом\\go-remerge\\internal\\arango\\client.go`)

																																																																																																																																																																											// Print information about each imported struct object

																																																																																																																																																																												for _, structObj := range importedStructObjects  {

																																																																																																																																																																															fmt.Printf("Imported struct object '%s' from package '%s' in file '%s'\n", structObj["struct_name"], structObj["package_name"], structObj["file_name"])

																																																																																																																																																																																}

																																																																																																																																																																															}

																																																																																																																																																																												}
																																																																																																																																																																							}
																																																																																																							}
																																																																														}
																																																																						}
																																																															}
																																																									}
																																																				}
																																																}
																																					}
																				}
												}
					}
