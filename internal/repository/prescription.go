package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/model"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// PrescriptionRepository is an interface that has all the function to be implemented inside health check repository
	PrescriptionRepository interface {
		Insert(ctx context.Context, req *model.PrescriptionRequest, phoneNumber string) error
		GetByPatientID(ctx context.Context, patientID string) ([]model.Prescription, error)
	}

	// PrescriptionRepositoryImpl is an app health check struct that consists of all the dependencies needed for perscription repository
	PrescriptionRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgxpool.Pool
	}
)

// NewPrescriptionRepository return new instances prescription repository
func NewPrescriptionRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *pgxpool.Pool) *PrescriptionRepositoryImpl {
	return &PrescriptionRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (pr *PrescriptionRepositoryImpl) Insert(ctx context.Context, req *model.PrescriptionRequest, phoneNumber string) error {
	qInsertMedication := `
		INSERT INTO medication (ref_id,identifier,code,code_display,form_code,form_value,amount,status,manufacturer,extension,batch) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id
	`
	qInsertMedicationIngredients := `
		INSERT INTO medication_ingredient (
			medication_id,code,display,is_active,strength_denominator,strength_numerator
		) VALUES %s
	`
	qCheckPatient := `
		SELECT id FROM patient WHERE ref_id = $1
	`
	qInsertPatient := `
		INSERT INTO patient (ref_id, name, phone_number) VALUES ($1, $2, $3) RETURNING id;
	`
	qInsertMedicationRequest := `
		INSERT INTO medication_request (medication_id,ref_id,status,patient_id,reason,intent,category,reported,encounter,requester,performer,recorder,note,insurance,course_of_therapy_type,dosage_instructions,dispense_request,substitution,raw_request) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id
	`
	qInsertMedicationRequestIdentifier := `
		INSERT INTO medication_request_identifier (
			medication_request_id,system,use,value
		) VALUES %s
	`

	tx, err := pr.DB.Begin(ctx)
	if err != nil {
		pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR begin TX", err)

		return err
	}

	// INSERT MEDICATION

	amountJsonData, err := json.Marshal(req.Medication.Amount)
	if err != nil {
		amountJsonData = []byte("{}")
	}

	batchJsonData, err := json.Marshal(req.Medication.Batch)
	if err != nil {
		batchJsonData = []byte("{}")
	}

	extJsonData, err := json.Marshal(req.Medication.Extension)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}

		pr.Logger.Error("PrescriptionRepositoryImpl.Insert json marshal ERROR", err)

		return err
	}

	var medicationID int
	row := tx.QueryRow(ctx, qInsertMedication,
		req.Medication.ID,
		req.Medication.Identifier[0].Value,
		req.Medication.Code.Coding[0].Code,
		req.Medication.Code.Coding[0].Display,
		req.Medication.Form.Coding[0].Code,
		req.Medication.Form.Coding[0].Display,
		string(amountJsonData),
		req.Medication.Status,
		req.Medication.Manufacturer.Reference,
		string(extJsonData),
		string(batchJsonData),
	)

	err = row.Scan(&medicationID)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}

		pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR Scan Insert Medication", err)

		return err
	}

	// INSERT BULK MEDICATION INGREDIENTS

	numberArgsPerRowMedicationIngredients := 6
	valueArgsMedicationIngredients := make([]interface{}, 0, numberArgsPerRowMedicationIngredients*len(req.Medication.Ingredient))

	for i := 0; i < len(req.Medication.Ingredient); i++ {
		valueArgsMedicationIngredients = append(valueArgsMedicationIngredients, medicationID, req.Medication.Ingredient[i].ItemCodeableConcept.Coding[0].Code, req.Medication.Ingredient[i].ItemCodeableConcept.Coding[0].Display, req.Medication.Ingredient[i].IsActive, fmt.Sprintf("%2.f %s", req.Medication.Ingredient[i].Strength.Denominator.Value, req.Medication.Ingredient[i].Strength.Denominator.Code), fmt.Sprintf("%2.f %s", req.Medication.Ingredient[i].Strength.Numerator.Value, req.Medication.Ingredient[i].Strength.Numerator.Code))
	}

	qInsertMedicationIngredients = helper.BulkInsert(qInsertMedicationIngredients, numberArgsPerRowMedicationIngredients, len(req.Medication.Ingredient))

	_, err = tx.Exec(ctx, qInsertMedicationIngredients, valueArgsMedicationIngredients...)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}

		pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR Exec Insert Bulk Medication Ingredients", err)

		return err
	}

	// INSERT PATIENT

	// check patient by ref id
	var (
		patientID        int
		trimPatientRefId string = strings.Replace(req.MedicationRequest.Subject.Reference, "Patient/", "", 1)
	)

	row = tx.QueryRow(ctx, qCheckPatient, trimPatientRefId)
	err = row.Scan(&patientID)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() { // exclude ErrNoRows first, for clean handling purposes
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}

		pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR Scan Check Patient", err)

		return err
	}

	// check error again, if ref id patient is not exists, then create patient
	if err != nil && err.Error() == pgx.ErrNoRows.Error() {
		row = tx.QueryRow(ctx, qInsertPatient, trimPatientRefId, req.MedicationRequest.Subject.Display, phoneNumber)

		err = row.Scan(&patientID)
		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

				return errRollback
			}

			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR Scan Insert Patient", err)

			return err
		}
	}

	// INSERT MEDICATION REQUEST

	noteJsonData, err := json.Marshal(req.MedicationRequest.Note)
	if err != nil {
		noteJsonData = []byte("{}")
	}

	insuranceJsonData, err := json.Marshal(req.MedicationRequest.Insurance)
	if err != nil {
		insuranceJsonData = []byte("{}")
	}

	dosageInstructionJsonData, err := json.Marshal(req.MedicationRequest.DosageInstruction)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}
		pr.Logger.Error("PrescriptionRepositoryImpl.Insert json marshal ERROR", err)

		return err
	}

	dispenseRequestJsonData, err := json.Marshal(req.MedicationRequest.DispenseRequest)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}
		pr.Logger.Error("PrescriptionRepositoryImpl.Insert json marshal ERROR", err)

		return err
	}

	substitutionJsonData, err := json.Marshal(req.MedicationRequest.Substitution)
	if err != nil {
		substitutionJsonData = []byte("{}")
	}

	rawRequestJsonData, err := json.Marshal(req.MedicationRequest)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}
		pr.Logger.Error("PrescriptionRepositoryImpl.Insert json marshal ERROR", err)

		return err
	}

	var medicationRequestID int
	row = tx.QueryRow(ctx, qInsertMedicationRequest,
		medicationID,
		req.MedicationRequest.ID,
		req.MedicationRequest.Status,
		patientID,
		req.MedicationRequest.ReasonCode[0].Coding[0].Display,
		req.MedicationRequest.Intent,
		req.MedicationRequest.Category[0].Coding[0].Code,
		false,
		req.MedicationRequest.Encounter.Reference,
		req.MedicationRequest.Requester.Reference,
		req.MedicationRequest.Performer.Reference,
		req.MedicationRequest.Recorder.Reference,
		string(noteJsonData),
		string(insuranceJsonData),
		req.MedicationRequest.CourseOfTherapyType.Coding[0].Code,
		string(dosageInstructionJsonData),
		string(dispenseRequestJsonData),
		string(substitutionJsonData),
		string(rawRequestJsonData),
	)

	err = row.Scan(&medicationRequestID)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}

		pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR Scan Insert Medication Request", err)

		return err
	}

	// INSERT BULK MEDICATION REQUEST IDENTIFIER
	numberArgsPerRowMedicationRequestIdentifiers := 4
	valueArgsMedicationRequestIdentifiers := make([]interface{}, 0, numberArgsPerRowMedicationRequestIdentifiers*len(req.MedicationRequest.Identifier))

	for i := 0; i < len(req.MedicationRequest.Identifier); i++ {
		valueArgsMedicationRequestIdentifiers = append(valueArgsMedicationRequestIdentifiers, medicationRequestID, req.MedicationRequest.Identifier[i].System, req.MedicationRequest.Identifier[i].Use, req.MedicationRequest.Identifier[i].Value)
	}

	qInsertMedicationRequestIdentifier = helper.BulkInsert(qInsertMedicationRequestIdentifier, numberArgsPerRowMedicationRequestIdentifiers, len(req.MedicationRequest.Identifier))

	_, err = tx.Exec(ctx, qInsertMedicationRequestIdentifier, valueArgsMedicationRequestIdentifiers...)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}

		pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR Exec Insert Bulk Medication Request Identifiers", err)

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return errRollback
		}

		pr.Logger.Error("PrescriptionRepositoryImpl.Insert ERROR commit TX", err)

		return err
	}

	return nil
}

func (pr *PrescriptionRepositoryImpl) GetByPatientID(ctx context.Context, patientID string) ([]model.Prescription, error) {
	q := `
		SELECT 
			m.id,
			m.code,
			m.code_display as display
		FROM
			medication m
		JOIN
			medication_request mr
		ON
			mr.medication_id = m.id
		JOIN
			patient p
		ON
			mr.patient_id = p.id
		WHERE 
			p.ref_id = $1
	`

	var prescriptions []model.Prescription
	rows, err := pr.DB.Query(ctx, q, patientID)
	if err != nil {
		pr.Logger.Error("PrescriptionRepositoryImpl.GetByPatientID Query ERROR", err)

		return []model.Prescription{}, err
	}

	for rows.Next() {
		prescription := model.Prescription{}

		err := rows.Scan(
			&prescription.ID,
			&prescription.Code,
			&prescription.Display,
		)
		if err != nil {
			pr.Logger.Error("PrescriptionRepositoryImpl.GetByPatientID rows Scan ERROR", err)
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, nil
}
