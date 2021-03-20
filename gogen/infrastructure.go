package gogen

import "fmt"

func createUtil(folderPath string) {
	CreateFolder("%s/infrastructure/util", folderPath)

	_ = WriteFileIfNotExist(
		"infrastructure/util/helper._go",
		fmt.Sprintf("%s/infrastructure/util/helper.go", folderPath),
		struct{}{},
	)
}

func createLog(folderPath string) {
	CreateFolder("%s/infrastructure/log", folderPath)

	_ = WriteFileIfNotExist(
		"infrastructure/log/contract._go",
		fmt.Sprintf("%s/infrastructure/log/contract.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"infrastructure/log/implementation._go",
		fmt.Sprintf("%s/infrastructure/log/implementation.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"infrastructure/log/public._go",
		fmt.Sprintf("%s/infrastructure/log/public.go", folderPath),
		struct{}{},
	)
}

func createAppError(folderPath string) {
	CreateFolder("%s/application/apperror", folderPath)

	_ = WriteFileIfNotExist(
		"application/apperror/error_function._go",
		fmt.Sprintf("%s/application/apperror/error_function.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/apperror/error_enum._go",
		fmt.Sprintf("%s/application/apperror/error_enum.go", folderPath),
		struct{}{},
	)
}

func createDomain(folderPath string) {
	CreateFolder("%s/domain/entity", folderPath)

	CreateFolder("%s/domain/repository", folderPath)

	CreateFolder("%s/domain/service", folderPath)

	_ = WriteFileIfNotExist(
		"domain/repository/repository._go",
		fmt.Sprintf("%s/domain/repository/repository.go", folderPath),
		struct{}{},
	)
}

func createDatabase(folderPath string) {
	CreateFolder("%s/domain/repository", folderPath)

	_ = WriteFileIfNotExist(
		"domain/repository/database._go",
		fmt.Sprintf("%s/domain/repository/database.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"domain/repository/transaction._go",
		fmt.Sprintf("%s/domain/repository/transaction.go", folderPath),
		struct{}{},
	)

}

func createConfig(folderPath string) {
	CreateFolder("%s/infrastructure/config", folderPath)

	_ = WriteFileIfNotExist(
		"infrastructure/config/config._go",
		fmt.Sprintf("%s/infrastructure/config/config.go", folderPath),
		struct{}{},
	)
}
